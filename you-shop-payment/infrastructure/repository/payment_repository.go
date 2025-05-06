package repository

import (
	"database/sql"

	"github.com/TechwizsonORG/payment-service/entity"
	"github.com/TechwizsonORG/payment-service/util"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(payment *entity.Payment) error {
	query := `
		INSERT INTO payment (id, user_id, order_id, amount, currency, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	payment.CreatedAt = util.GetCurrentUtcTime(7)
	payment.UpdatedAt = util.GetCurrentUtcTime(7)
	_, err := r.db.Exec(
		query,
		payment.Id,
		payment.UserId,
		payment.OrderId,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	return err
}

func (r *PaymentRepository) CreateTransaction(transaction *entity.Transaction) error {
	query := `
		INSERT INTO transaction (id, payment_id, amount, currency, type, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(
		query,
		transaction.Id,
		transaction.PaymentId,
		transaction.Amount,
		transaction.Currency,
		transaction.Type,
		transaction.Status,
		util.GetCurrentUtcTime(7),
	)
	return err
}

func (r *PaymentRepository) UpdatePaymentStatus(paymentID uuid.UUID, paymentStatus entity.PaymentStatus, transactionStatus entity.TransactionStatus) error {
	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Update payment status
	query := `
        UPDATE payment
        SET status = $1, updated_at = $2
        WHERE id = $3
    `
	_, err = tx.Exec(query, paymentStatus, util.GetCurrentUtcTime(7), paymentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create a new transaction
	transaction := &entity.Transaction{
		AuditEntity: entity.AuditEntity{
			Id:        uuid.New(),
			CreatedAt: util.GetCurrentUtcTime(7),
			UpdatedAt: util.GetCurrentUtcTime(7),
		},
		PaymentId: paymentID,
		Status:    transactionStatus,
	}

	query = `
        INSERT INTO transaction (id, payment_id, amount, status, created_at, updated_at)
        VALUES ($1, $2, (SELECT p.amount FROM payment p where p.id = $3), $4, $5, $6)
    `
	_, err = tx.Exec(query, transaction.Id, transaction.PaymentId, transaction.PaymentId, transaction.Status, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *PaymentRepository) GetUserPayments(userId uuid.UUID, pageSize, offset int) ([]*entity.Payment, error) {
	query := `
		SELECT id, user_id, order_id, amount, currency, status, created_at, updated_at
		FROM payment WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := p.db.Query(query, userId, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payemnts := []*entity.Payment{}
	for rows.Next() {
		payment := &entity.Payment{}
		err := rows.Scan(
			&payment.Id,
			&payment.UserId,
			&payment.OrderId,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payemnts = append(payemnts, payment)
	}
	return payemnts, nil
}

func (p *PaymentRepository) GetPaymentById(paymentId uuid.UUID) (*entity.Payment, error) {
	query := `
		SELECT id, user_id, order_id, amount, currency, status, created_at, updated_at
		FROM payment WHERE id = $1
	`
	row := p.db.QueryRow(query, paymentId)
	payment := &entity.Payment{}
	err := row.Scan(
		&payment.Id,
		&payment.UserId,
		&payment.OrderId,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (p *PaymentRepository) GetPaymensByOrderIds(orderIds []uuid.UUID) ([]*entity.Payment, error) {
	query := `
        SELECT id, user_id, order_id, amount, currency, status, created_at, updated_at
        FROM payment
        WHERE order_id = ANY($1)
    `
	rows, err := p.db.Query(query, pq.Array(orderIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []*entity.Payment{}
	for rows.Next() {
		payment := &entity.Payment{}
		err := rows.Scan(
			&payment.Id,
			&payment.UserId,
			&payment.OrderId,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (p *PaymentRepository) GetPaymentByOrderId(orderId uuid.UUID) (*entity.Payment, error) {
	query := `
		SELECT id, user_id, order_id, amount, currency, status, created_at, updated_at
		FROM payment WHERE order_id = $1 AND status = 1
	`
	row := p.db.QueryRow(query, orderId)
	payment := &entity.Payment{}
	err := row.Scan(
		&payment.Id,
		&payment.UserId,
		&payment.OrderId,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No payment found for the given order ID
		}
		return nil, err
	}
	return payment, nil
}
