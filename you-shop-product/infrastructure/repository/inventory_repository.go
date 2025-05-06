package repository

import (
	"database/sql"

	"github.com/TechwizsonORG/product-service/entity"
	"github.com/TechwizsonORG/product-service/usecase/inventory/model"
	"github.com/TechwizsonORG/product-service/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type InventoryRepository struct {
	db  *sql.DB
	log zerolog.Logger
}

func NewInventoryRepository(db *sql.DB, log zerolog.Logger) *InventoryRepository {
	logger := log.
		With().
		Str("repository", "inventory").
		Logger()
	return &InventoryRepository{db: db, log: logger}

}
func (p *InventoryRepository) GetQuantity(productId, sizeId, colorId uuid.UUID) (int, error) {
	query := `
		SELECT
			COALESCE(i.quantity, 0) AS quantity
		FROM inventory i
		JOIN product p ON i.product_id = p.id
		WHERE i.product_id = $1 AND i.size_id = $2 AND i.color_id = $3 AND p.deleted_at IS NULL AND p.status = 1
	`
	row := p.db.QueryRow(query, productId, sizeId, colorId)
	var quantity int
	scanErr := row.Scan(&quantity)
	if scanErr != nil {
		p.log.Error().Err(scanErr).Msg("")
		return 0, scanErr
	}
	return quantity, nil
}

func (p *InventoryRepository) AddInventories(createInventories []model.CreateInventory) error {
	query := `
		INSERT INTO inventory (product_id, size_id, color_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	for _, inventory := range createInventories {
		stmt, err := tx.Prepare(query)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(inventory.ProductId, inventory.SizeId, inventory.ColorId, inventory.Quantity, util.GetCurrentUtcTime(7), util.GetCurrentUtcTime(7))
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (i *InventoryRepository) GetInventory(productId uuid.UUID, colorId uuid.UUID, sizeId uuid.UUID) (inventory *entity.Inventory, e error) {
	query := `
		SELECT
			i.product_id,
			i.size_id,
			i.color_id,
			i.quantity
		FROM inventory i
		WHERE i.color_id = $1
			AND i.product_id = $2
			AND i.size_id = $3
	`
	row := i.db.QueryRow(query, colorId, productId, sizeId)
	if row.Err() != nil {
		return nil, row.Err()
	}
	inventory = &entity.Inventory{}
	scanErr := row.Scan(&inventory.ProductId, &inventory.SizeId, &inventory.ColorId, &inventory.Quantity)
	if scanErr != nil {
		if scanErr.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, scanErr
	}
	return inventory, nil
}
func (i *InventoryRepository) UpdateInventory(updateInventory *entity.Inventory) error {
	query := `
		UPDATE inventory
		SET
			quantity = $1
		WHERE color_id = $2
			AND product_id = $3
			AND size_id = $4
	`
	stmt, err := i.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(updateInventory.Quantity, updateInventory.ColorId, updateInventory.ProductId, updateInventory.SizeId)
	if err != nil {
		return err
	}
	return nil
}
