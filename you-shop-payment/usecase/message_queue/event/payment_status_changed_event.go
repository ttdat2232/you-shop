package event

import (
	"github.com/TechwizsonORG/payment-service/entity"
	"github.com/google/uuid"
)

type PaymentStatusChangedEvent struct {
	PaymentId uuid.UUID            `json:"payementId"`
	Status    entity.PaymentStatus `json:"status"`
	OrderId   uuid.UUID            `json:"orderId"`
}
