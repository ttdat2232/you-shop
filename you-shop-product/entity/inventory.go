package entity

import "github.com/google/uuid"

type Inventory struct {
	ColorId   uuid.UUID
	ProductId uuid.UUID
	SizeId    uuid.UUID
	Quantity  int
}
