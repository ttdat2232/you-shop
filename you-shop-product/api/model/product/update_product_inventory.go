package product

import "github.com/google/uuid"

type UpdateProductInventoryRequest struct {
	SizeId   uuid.UUID
	ColorId  uuid.UUID
	Price    float64
	Quantity int
}
