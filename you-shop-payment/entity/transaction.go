package entity

import "github.com/google/uuid"

type TransactionStatus int8

const (
	TransactionSuccess TransactionStatus = iota + 1
	TransciontFailed
)

type Transaction struct {
	AuditEntity
	PaymentId uuid.UUID
	Amount    float64
	Currency  string
	Type      string            // e.g., DEBIT, CREDIT
	Status    TransactionStatus // e.g., SUCCESS, FAILED
}
