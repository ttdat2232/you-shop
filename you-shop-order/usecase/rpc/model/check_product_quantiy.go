package model

import "github.com/google/uuid"

type CheckProductQuantity struct {
	ProductId       uuid.UUID `json:"productId"`
	SizeId          uuid.UUID `json:"sizeId"`
	ColorId         uuid.UUID `json:"colorId"`
	RequireQuantity int       `json:"requireQuantity"`
}
