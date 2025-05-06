package model

import "github.com/google/uuid"

type TotalPriceRequest struct {
	Items []*OrderItem `json:"items"`
}

type OrderItem struct {
	ProductId uuid.UUID `json:"productId"`
	SizeId    uuid.UUID `json:"sizeId"`
	ColorId   uuid.UUID `json:"colorId"`
	Quantity  int       `json:"quantity"`
}
