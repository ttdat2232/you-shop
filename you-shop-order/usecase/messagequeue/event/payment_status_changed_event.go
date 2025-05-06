package event

import "github.com/google/uuid"

type PaymentStatus int8

const (
	PaymentPending PaymentStatus = iota + 1
	PaymentSuccess
	PaymentFailed
)

type PaymentStatusChangedEvent struct {
	PaymentId uuid.UUID     `json:"payementId"`
	Status    PaymentStatus `json:"status"`
	OrderId   uuid.UUID     `json:"orderId"`
}
