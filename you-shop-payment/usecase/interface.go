package usecase

import (
	"github.com/TechwizsonORG/payment-service/entity"
	"github.com/TechwizsonORG/payment-service/err"
	"github.com/TechwizsonORG/payment-service/usecase/model"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	CreatePayment(payment *entity.Payment) error
	CreateTransaction(transaction *entity.Transaction) error
	UpdatePaymentStatus(paymentID uuid.UUID, paymentStatus entity.PaymentStatus, transactionStatus entity.TransactionStatus) error
	GetUserPayments(userId uuid.UUID, pageSize, offset int) ([]*entity.Payment, error)
	GetPaymentById(paymentId uuid.UUID) (*entity.Payment, error)
	GetPaymensByOrderIds(orderIds []uuid.UUID) ([]*entity.Payment, error)
	GetPaymentByOrderId(orderId uuid.UUID) (*entity.Payment, error)
}

type Service interface {
	CreatePayment(order model.CreateOrder) (*entity.Payment, err.AppError)
	ProcessPaymentCallback(paymentId uuid.UUID, transactionStatus entity.TransactionStatus) err.AppError
	GetUserPayments(userId uuid.UUID, pageSize, pageNumber int) []*entity.Payment
	GetPaymentsByOrderIds(orderIds []uuid.UUID) []*entity.Payment
	GetPaymentByOrderId(orderId uuid.UUID) (*entity.Payment, err.AppError)
}
