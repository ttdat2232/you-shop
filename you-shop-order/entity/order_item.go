package entity

import "github.com/google/uuid"

type OrderItem struct {
	Quantity  int
	ProductId uuid.UUID
	ColorId   uuid.UUID
	SizeId    uuid.UUID
	OrderId   uuid.UUID
	Price     float64
	PriceId   uuid.UUID
}
