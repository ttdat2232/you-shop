package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrder struct {
	UserId      uuid.UUID          `json:"-"`
	Items       []*CreateOrderItem `json:"items"`
	Description string             `json:"description"`
}

type CreateOrderItem struct {
	ProductId uuid.UUID `json:"productId"`
	SizeId    uuid.UUID `json:"sizeId"`
	ColorId   uuid.UUID `json:"colorId"`
	Quantity  int       `json:"quantity"`
}

type OrderResponse struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	TotalPrice  float64   `json:"totalPrice"`
	CreatedAt   time.Time `json:"createdAt"`
}
