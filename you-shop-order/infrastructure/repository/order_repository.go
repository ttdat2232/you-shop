package repository

import (
	"database/sql"

	"github.com/TechwizsonORG/order-service/entity"
	"github.com/TechwizsonORG/order-service/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type OrderRepository struct {
	logger zerolog.Logger
	db     *sql.DB
}

func NewOrderRepository(db *sql.DB, logger zerolog.Logger) *OrderRepository {
	repoLogger := logger.With().Str("Infrastructure", "OrderRepository").Logger()
	return &OrderRepository{
		db:     db,
		logger: repoLogger,
	}
}

func (o *OrderRepository) GetOrders(pageIndex int, pageSize int) ([]entity.Order, error) {
	return nil, nil
}
func (o *OrderRepository) GetOrderById(orderId uuid.UUID) (*entity.Order, error) {
	query := `
		SELECT
			o.id,
			o.description,
			o.total_price,
			o.status,
			o.owner_id,
			o.created_at,
			o.updated_at,
			oi.quantity,
			oi.product_id,
			oi.color_id,
			oi.size_id,
			oi.price,
			oi.price_id
		FROM "order" o
		JOIN order_item oi ON o.id = oi.order_id
		WHERE o.is_deleted = false AND o.id = $1
	`
	rows, err := o.db.Query(query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	order := &entity.Order{}
	for rows.Next() {
		var item entity.OrderItem
		err := rows.Scan(
			&order.Id,
			&order.Description,
			&order.TotalPrice,
			&order.Status,
			&order.OwnerId,
			&order.CreatedAt,
			&order.UpdatedAt,
			&item.Quantity,
			&item.ProductId,
			&item.ColorId,
			&item.SizeId,
			&item.Price,
			&item.PriceId,
		)
		order.Items = append(order.Items, &item)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}
func (o *OrderRepository) GetOrdersByUserId(userId uuid.UUID, pageIndex int, pageSize int) ([]entity.Order, error) {
	query := `
		SELECT
			o.id,
			o.description,
			o.total_price,
			o.status,
			o.created_at
		FROM "order" o
		WHERE o.is_deleted = false AND o.owner_id = $1
		ORDER BY o.created_at DESC
		OFFSET $2
		LIMIT $3
	`
	offset := pageIndex * pageSize
	rows, err := o.db.Query(query, userId, offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []entity.Order{}
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.Id, &order.Description, &order.TotalPrice, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}
	return result, nil
}
func (o *OrderRepository) IsOrderExistById(orderId uuid.UUID) (bool, error) {
	query := `
		SELECT o.id FROM "order" o
		WHERE o.id = $1 AND o.is_deleted = false
	`
	rows, err := o.db.Query(query, orderId)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	counter := 0
	for rows.Next() {
		counter++
	}
	return counter > 0, nil
}

func (o *OrderRepository) CreateOrder(order *entity.Order) error {
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	query := `
		INSERT INTO "order" (id, description, total_price, status, owner_id, order_code)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(order.Id, order.Description, order.TotalPrice, order.Status, order.OwnerId, order.OrderCode)
	if err != nil {
		tx.Rollback()
		return err
	}
	itemQuery := `
		INSERT INTO order_item (quantity, product_id, color_id, size_id, order_id, price, price_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	itemStmt, err := tx.Prepare(itemQuery)
	if err != nil {
		return err
	}
	defer itemStmt.Close()
	for _, item := range order.Items {
		_, err = itemStmt.Exec(item.Quantity, item.ProductId, item.ColorId, item.SizeId, order.Id, item.Price, item.PriceId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderRepository) UpdateOrder(entity *entity.Order) error {
	entity.UpdatedAt = util.GetCurrentUtcTime(7)
	query := `
		UPDATE "order"
		SET
			description = $1,
			status = $2,
			order_code = $3,
			total_price = $4,
			owner_id = $5,
			updated_at = $6,
			is_deleted = $7,
			deleted_at = $8
		WHERE id = $9
	`
	stmt, err := o.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(entity.Description, entity.Status, entity.OrderCode, entity.TotalPrice, entity.OwnerId, entity.UpdatedAt, entity.IsDeleted, entity.DeletedAt, entity.Id)
	if err != nil {
		return err
	}
	return nil
}
func (o *OrderRepository) DeleteOrder(*entity.Order) error {
	return nil
}
