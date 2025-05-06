package model

import "github.com/google/uuid"

type CreateInventory struct {
	ProductId uuid.UUID
	ColorId   uuid.UUID
	SizeId    uuid.UUID
	Price     float64
	Quantity  int
}
