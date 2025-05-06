package event

import (
	"github.com/google/uuid"
)

type ProductCreatedEvent struct {
	ProductId uuid.UUID `json:"productId"`
	Price     float64   `json:"price"`
}
