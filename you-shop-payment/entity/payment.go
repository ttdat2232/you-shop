package entity

import (
	"github.com/google/uuid"
)

type PaymentStatus int8

const (
	PaymentPending PaymentStatus = iota + 1
	PaymentSuccess
	PaymentFailed
)

type Payment struct {
	AuditEntity
	OrderId         uuid.UUID
	UserId          uuid.UUID
	Amount          float64
	Currency        string
	Status          PaymentStatus
	PaymentMethodId uuid.UUID
}

func NewPayment(amount float64, orderId, userId uuid.UUID) *Payment {
	return &Payment{
		UserId:  userId,
		OrderId: orderId,
		Amount:  amount,
		Status:  PaymentPending,
		AuditEntity: AuditEntity{
			Id: uuid.New(),
		},
	}
}

type PaymentMethod struct {
	AuditEntity
	Type      string // e.g., CREDIT_CARD, BANK_ACCOUNT
	Details   string // Encrypted JSON
	IsDefault bool
}
