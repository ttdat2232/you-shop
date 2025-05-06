package model

import "github.com/google/uuid"

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
