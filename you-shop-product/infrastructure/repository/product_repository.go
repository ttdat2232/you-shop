package repository

import (
	"database/sql"
	"time"

	"github.com/TechwizsonORG/product-service/entity"
	appErr "github.com/TechwizsonORG/product-service/err"
	"github.com/TechwizsonORG/product-service/util"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

type ProductRepository struct {
	db  *sql.DB
	log zerolog.Logger
}

func NewProductRepository(db *sql.DB, log zerolog.Logger) *ProductRepository {
	logger := log.
		With().
		Str("repository", "product").
		Logger()
	return &ProductRepository{db: db, log: logger}

}

func (p *ProductRepository) Search(query string) ([]entity.Product, appErr.ApplicationError) {
	queryString := `
		SELECT
			p.id::uuid,
			p.status,
			p.name,
			p.description,
			p.sku,
			p.created_at,
			p.updated_at,
			p.thumbnail
		FROM product p
		WHERE p.name LIKE $1 AND p.deleted_at IS NULL
	`
	var products []entity.Product
	rows, err := p.db.Query(queryString, "%"+query+"%")
	if err != nil {
		p.log.Error().Err(err).Msg("Errored occurred")
		return nil, appErr.CommonError()
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		var status entity.ProductStatus
		var name sql.NullString
		var description sql.NullString
		var sku sql.NullString
		var createdAt time.Time
		var updatedAt time.Time
		var thumbnail sql.NullString
		err := rows.Scan(&id, &status, &name, &description, &sku, &createdAt, &updatedAt, &thumbnail)
		if err != nil {
			p.log.Err(err).Msg("Error occurred")
			return nil, appErr.CommonError()
		}
		product := entity.Product{
			Id:          id,
			Status:      status,
			Name:        name.String,
			Description: description.String,
			Sku:         sku.String,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Thumbnail:   thumbnail.String,
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductRepository) Count() (int, appErr.ApplicationError) {
	queryString := `
		SELECT COUNT(*) FROM product p
		WHERE p.deleted_at IS NULL
	`
	row := p.db.QueryRow(queryString)
	var count int
	err := row.Scan(&count)
	if err != nil {
		p.log.Error().Err(err).Msg("")
		return 0, nil
	}

	return count, nil
}

func (p *ProductRepository) List(page int, pageSize int) ([]entity.Product, appErr.ApplicationError) {
	queryString := `
		SELECT
			p.id::uuid,
			p.status,
			p.name,
			p.description,
			p.sku,
			p.created_at,
			p.updated_at,
			p.thumbnail
		FROM product p
		WHERE p.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT $1
		OFFSET $2
	`
	offset := (page - 1) * pageSize
	rows, err := p.db.Query(queryString, pageSize, offset)
	if err != nil {
		p.log.Error().Err(err).Msg("")
		return nil, appErr.CommonError()
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var id uuid.UUID
		var status entity.ProductStatus
		var name sql.NullString
		var description sql.NullString
		var sku sql.NullString
		var createdAt time.Time
		var updatedAt time.Time
		var thumbnail sql.NullString
		err := rows.Scan(&id, &status, &name, &description, &sku, &createdAt, &updatedAt, &thumbnail)
		if err != nil {
			p.log.Err(err).Msg("")
			return nil, appErr.CommonError()
		}
		product := entity.Product{
			Id:          id,
			Status:      status,
			Name:        name.String,
			Description: description.String,
			Sku:         sku.String,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Thumbnail:   thumbnail.String,
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *ProductRepository) Get(id uuid.UUID) (*entity.Product, appErr.ApplicationError) {
	query := `
		SELECT
			p.id::uuid,
			p.status,
			p.name, p.description,
			p.sku,
			p.created_at,
			p.updated_at,
			p.thumbnail
		FROM product p
		WHERE id = $1 AND p.deleted_at IS NULL
	`
	row, err := p.db.Query(query, id.String())
	if err != nil {
		p.log.Error().Err(err).Msg("")
		return nil, appErr.CommonError()
	}
	defer row.Close()

	var status entity.ProductStatus
	var name sql.NullString
	var description sql.NullString
	var sku sql.NullString
	var createdAt time.Time
	var updatedAt time.Time
	var thumbnail sql.NullString
	for row.Next() {
		err = row.Scan(&id, &status, &name, &description, &sku, &createdAt, &updatedAt, &thumbnail)
		if err != nil {
			p.log.Error().Err(err).Msg("")
			return nil, appErr.CommonError()
		}

		product := entity.Product{
			Id:          id,
			Status:      status,
			Name:        name.String,
			Description: description.String,
			Sku:         sku.String,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Thumbnail:   thumbnail.String,
		}
		return &product, nil
	}
	return nil, appErr.NotFoundProductErrorWithId(id.String())
}

func (p *ProductRepository) Create(product entity.Product) (entity.Product, appErr.ApplicationError) {
	query := `
		INSERT INTO product (id, name, description, sku, status, created_at, thumbnail, user_manual)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	stmt, prepareErr := p.db.Prepare(query)
	if prepareErr != nil {
		p.log.Error().Err(prepareErr).Msg("Error occurred")
		return entity.Product{}, appErr.NewProductError(500, "Error occurred", "Error occurred", nil)
	}
	_, err := stmt.Exec(product.Id, product.Name, product.Description, product.Sku, product.Status, product.CreatedAt, product.Thumbnail, product.UserManual)

	if err != nil {
		p.log.Error().Err(err).Msg("")
		return entity.Product{}, appErr.CommonError()
	}
	return product, nil
}

func (p *ProductRepository) Update(product entity.Product) (entity.Product, appErr.ApplicationError) {
	query := `
		UPDATE product
		SET
			name = $1,
			description = $2,
			sku = $3,
			status = $4,
			user_manual = $5,
			updated_at = $6
		WHERE id = $7
	`
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return entity.Product{}, appErr.CommonError()
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Description, product.Sku, product.Status, product.UserManual, util.GetCurrentUtcTime(7), product.Id)
	if err != nil {
		p.log.Error().Err(err).Msg("")
		return entity.Product{}, appErr.CommonError()
	}
	return product, nil
}

func (p *ProductRepository) Delete(id uuid.UUID) appErr.ApplicationError {
	query := `
		UPDATE product SET deleted_at = $1 WHERE id = $2
	`
	stmt, err := p.db.Prepare(query)
	if err != nil {
		p.log.Error().Err(err).Msg("")
		return appErr.CommonError()
	}
	defer stmt.Close()

	_, err = stmt.Exec(util.GetCurrentUtcTime(7), id)
	if err != nil {
		p.log.Error().Err(err).Msg("")
		return appErr.CommonError()
	}
	return nil
}

func (p *ProductRepository) IsSkuAlreadyExisted(sku string) (bool, appErr.ApplicationError) {
	queryStr := `
		SELECT COUNT(p.id) FROM product p
		WHERE p.sku like $1
	`
	row := p.db.QueryRow(queryStr, sku)
	var count int
	scanErr := row.Scan(&count)
	if scanErr != nil {
		p.log.Error().Err(scanErr).Msg("")
		return false, nil
	}

	if count == 0 || count > 1 {
		return false, nil
	}
	return true, nil
}

func (p *ProductRepository) IsIdExisted(id uuid.UUID) (bool, appErr.ApplicationError) {
	queryStr := `
		SELECT COUNT(p.id) FROM product p
		WHERE p.id = $1 AND p.deleted_at IS NULL
	`
	row := p.db.QueryRow(queryStr, id)
	var count int
	scanErr := row.Scan(&count)
	if scanErr != nil {
		p.log.Error().Err(scanErr).Msg("")
		return false, nil
	}

	if count == 0 || count > 1 {
		return false, nil
	}
	return true, nil

}

func (p *ProductRepository) GetByIds(productIds []uuid.UUID) ([]entity.Product, appErr.ApplicationError) {
	query := `
		SELECT
			p.id,
			COALESCE(name, 'Unknown') AS name
			COALESCE(thumbnail, '') AS thumbnail
		FROM
			UNNEST ($1::uuid[]) AS target_id(id)
		LEFT JOIN
			product p ON p.id = target_id.id
		WHERE p.status = 1 AND p.deleted_at IS NULL
	`

	rows, queryErr := p.db.Query(query, pq.Array(productIds))
	if queryErr != nil {
		p.log.Error().Err(queryErr).Msg("")
		return []entity.Product{}, nil
	}
	defer rows.Close()
	result := []entity.Product{}
	for rows.Next() {
		var id uuid.UUID
		var name string
		var thumbnail string
		scanErr := rows.Scan(&id, &name, &thumbnail)
		if scanErr != nil {
			p.log.Error().Err(scanErr).Msg("")
			return []entity.Product{}, nil
		}
		product := entity.Product{
			Id:        id,
			Name:      name,
			Thumbnail: thumbnail,
		}
		result = append(result, product)
	}
	return result, nil
}

func (s *ProductRepository) AddProductImages(productImages []entity.ProductImage) appErr.ApplicationError {
	query := `
		INSERT INTO product_image (id, product_id, color_id, image_url, is_primary, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	tx, err := s.db.Begin()
	if err != nil {
		s.log.Error().Err(err).Msg("")
		return appErr.CommonError()
	}
	for _, productImage := range productImages {
		stmt, err := tx.Prepare(query)
		if err != nil {
			s.log.Error().Err(err).Msg("")
			tx.Rollback()
			return appErr.CommonError()
		}
		defer stmt.Close()
		if _, err := stmt.Exec(productImage.Id, productImage.ProductId, productImage.ColorId, productImage.ImageUrl, productImage.IsPrimary, productImage.IsPublic, util.GetCurrentUtcTime(7), util.GetCurrentUtcTime(7)); err != nil {
			s.log.Error().Err(err).Msg("")
			tx.Rollback()
			return appErr.CommonError()
		}

	}
	tx.Commit()
	return nil
}
