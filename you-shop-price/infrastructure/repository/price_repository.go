package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechwizsonORG/price-service/entity"
	"github.com/TechwizsonORG/price-service/usecase/rpc/model"
	"github.com/TechwizsonORG/price-service/util"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

type PriceRepository struct {
	logger zerolog.Logger
	db     *sql.DB
}

func NewPriceRepository(logger zerolog.Logger, db *sql.DB) *PriceRepository {
	logger = logger.With().Str("repository", "price").Logger()
	return &PriceRepository{
		logger: logger,
		db:     db,
	}
}

func (p *PriceRepository) GetCurrentPrices(items []*model.OrderItem) ([]entity.Price, error) {
	keys := make([]interface{}, 0, len(items))
	for _, item := range items {
		keys = append(keys, item.ProductId, item.ColorId, item.SizeId)
	}

	queryBuff := strings.Builder{}
	queryBuff.WriteString(`
		SELECT
			p.id,
			COALESCE(p.amount, 0) AS amount,
			p.color_id,
			p.product_id,
			p.size_id
		FROM price p
		INNER JOIN price_list pl ON p.price_list_id = pl.id
		WHERE p.deleted_at IS NULL
			AND p.is_active = true
			AND pl.currency = 1
			AND NOW() BETWEEN p.valid_from AND COALESCE(p.valid_to, 'infinity'::timestamptz)
			AND (p.product_id, p.color_id, p.size_id) IN (
	`)
	for i := range items {
		if i > 0 {
			queryBuff.WriteString(", ")
		}

		queryBuff.WriteString(fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3))
	}
	queryBuff.WriteRune(')')
	query := queryBuff.String()
	p.logger.Debug().Msgf("Get current price query: %s", query)
	rows, queryErr := p.db.Query(query, keys...)

	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	result := make([]entity.Price, 0, len(items))
	for rows.Next() {
		var price entity.Price
		scanErr := rows.Scan(&price.Id, &price.Amount, &price.ColorId, &price.ProductId, &price.SizeId)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, price)
	}
	return result, nil
}

func (p *PriceRepository) GetDefaultPriceList() *entity.PriceList {
	queryStr := `
		SELECT pl.id, pl.description, pl.currency FROM price_list pl
		WHERE pl.currency = 1
		ORDER BY pl.created_at DESC
		LIMIT 1
	`
	row := p.db.QueryRow(queryStr)

	var priceList entity.PriceList
	row.Scan(&priceList.Id, &priceList.Description, &priceList.Currency)
	return &priceList
}

func (p *PriceRepository) AddNewPrice(price entity.Price) (*entity.Price, error) {
	err := p.AddNewPrices([]*entity.Price{&price})
	if err != nil {
		return nil, err
	}
	return &price, nil
}
func (p *PriceRepository) AddNewPriceList(priceList entity.PriceList) (*entity.PriceList, error) {
	queryStr := `
		INSERT INTO price_list (id, description, currency)
		VALUES ($1, $2, $3)
	`

	stmt, err := p.db.Prepare(queryStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(priceList.Id, priceList.Description, priceList.Currency)
	if err != nil {
		return nil, err
	}
	return &priceList, nil
}

func (p *PriceRepository) UpdatePrice(productId, colorId, sizeId uuid.UUID, price float64) (bool, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return false, err
	}

	priceHistoryInsertQuery := `
		INSERT INTO price_history (id, product_id, old_price_id, previous_price, new_price, change_reason)
		VALUES 
		(
			generate_uuid_v4(), 
			$1,
			(SELECT p.id FROM price p WHERE p.product_id = $2 AND p.color_id = $3 AND p.size_id = $4 AND p.is_active = true LIMIT 1),
			(SELECT p.amount FROM price p WHERE p.product_id = $5 AND p.color_id = $6 AND p.size_id = $7 AND p.is_active = true LIMIT 1),
			$8,
			'price updated'
		)
	`
	priceHistoryStmt, err := tx.Prepare(priceHistoryInsertQuery)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	defer priceHistoryStmt.Close()
	_, err = priceHistoryStmt.Exec(uuid.New(), productId, colorId, sizeId, productId, colorId, sizeId, price)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	priceUpdateQuery := `
		UPDATE price 
		SET valid_to = NOW(),
			is_active = false
		WHERE product_id = $1 AND color_id = $2 AND size_id = $3 AND is_active = true
	`
	priceUpdateStmt, err := tx.Prepare(priceUpdateQuery)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	defer priceUpdateStmt.Close()
	_, err = priceUpdateStmt.Exec(productId, colorId, sizeId)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	priceInsertQuery := `
		INSERT INTO price (id, valid_from, amount, product_id, color_id, size_id, is_active, price_list_id)
		VALUES (generate_uuid_v4(), $1, $2, $3, $4, $5, true, $6)
	`
	priceInsertStmt, err := tx.Prepare(priceInsertQuery)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	defer priceInsertStmt.Close()
	_, err = priceInsertStmt.Exec(util.GetCurrentUtcTime(7), price, productId, colorId, sizeId, p.GetDefaultPriceList().Id)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *PriceRepository) AddNewPrices(prices []*entity.Price) error {
	query := `
		INSERT INTO price (id, created_at, updated_at, amount, product_id, color_id, size_id, is_active, valid_from, price_list_id)	
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	priceList := p.GetDefaultPriceList()
	for _, price := range prices {
		stmt, err := tx.Prepare(query)
		if err != nil {
			tx.Rollback()
			return err
		}
		if _, err = stmt.Exec(price.Id, price.CreatedAt, price.UpdatedAt, price.Amount, price.ProductId, price.ColorId, price.SizeId, price.IsActive, price.ValidFrom, priceList.Id); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (p *PriceRepository) GetPrice(productId uuid.UUID, colorId uuid.UUID, sizeId uuid.UUID) (*entity.Price, error) {
	query := `
		SELECT id, created_at, updated_at, amount, product_id, color_id, size_id, is_active, valid_from, valid_to, price_list_id
		FROM price
		WHERE product_id = $1 AND color_id = $2 AND size_id = $3 AND is_active = true
		LIMIT 1
	`
	row := p.db.QueryRow(query, productId, colorId, sizeId)
	var price entity.Price
	var validTo sql.NullTime
	var validFrom sql.NullTime
	err := row.Scan(&price.Id, &price.CreatedAt, &price.UpdatedAt, &price.Amount, &price.ProductId, &price.ColorId, &price.SizeId, &price.IsActive, &validFrom, &validTo, &price.PriceListId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	price.ValidFrom = validFrom.Time
	price.ValidTo = validTo.Time
	return &price, nil
}

func (p *PriceRepository) GetCurrentPricesByProductIds(ids []uuid.UUID) ([]entity.Price, error) {
	query := `
		WITH ranked_prices AS (
			SELECT
				p.id,
				p.amount,
				p.product_id
				ROW_VERSION() OVER (PARTITION BY p.product_id ORDER BY p.amount) as rn
			FROM
        		price p
    		JOIN price_list pl ON p.price_list_id = pl.id
    		WHERE
        		p.is_active = true
        		AND p.deleted_at IS NULL
        		AND pl.currency = 1
        		AND NOW() BETWEEN p.valid_from AND COALESCE(p.valid_to, 'infinity'::timestamptz)
		)
		SELECT
			rp.id,
			COALESCE(p.amount, 0),
			target_id.id	
		FROM
			unnest($1::uuid[]) AS target_id(id)
		LEFT JOIN ranked_prices rp ON target_id.id = rp.product_id
		WHERE rp.rn = 1
	`
	rows, err := p.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []entity.Price{}
	for rows.Next() {
		var price entity.Price
		scanErr := rows.Scan(&price.Id, &price.Amount, &price.ProductId)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, price)
	}
	return result, nil
}
