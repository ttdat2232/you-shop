package product

import "github.com/google/uuid"

type CreateProductInventory struct {
	Inventories []Inventory
}

type Inventory struct {
	SizeId   uuid.UUID
	ColorId  uuid.UUID
	Price    float64
	Quantity int
}
