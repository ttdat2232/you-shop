package model

import "github.com/google/uuid"

type TotalPriceResponse struct {
	TotalPrice float64 `json:"totalPrice"`
	Items      []*Item `json:"items"`
}

type Item struct {
	Amount    float64   `json:"amount"`
	Quantity  int       `json:"quantity"`
	ProductId uuid.UUID `json:"productId"`
	ColorId   uuid.UUID `json:"colorId"`
	SizeId    uuid.UUID `json:"sizeId"`
	PriceId   uuid.UUID `json:"priceId"`
}
