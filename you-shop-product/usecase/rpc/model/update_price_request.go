package model

import "github.com/google/uuid"

type UpdatePriceRequest struct {
	ProductId uuid.UUID
	ColorId   uuid.UUID
	SizeId    uuid.UUID
	Price     float64
}

type UpdatePriceResponse struct {
	IsUpdated bool
}
