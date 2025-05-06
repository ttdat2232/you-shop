package model

import (
	"fmt"

	"github.com/TechwizsonORG/price-service/entity"
	"github.com/google/uuid"
)

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

func From(totalPrice float64, prices []entity.Price, orderItems []*OrderItem) *TotalPriceResponse {
	quantityMap := make(map[string]int, len(orderItems))
	for _, item := range orderItems {
		key := fmt.Sprintf("%s-%s-%s", item.ProductId, item.ColorId, item.SizeId)
		quantityMap[key] = item.Quantity
	}
	items := []*Item{}
	for _, price := range prices {
		key := fmt.Sprintf("%s-%s-%s", price.ProductId, price.ColorId, price.SizeId)
		items = append(items, &Item{
			Amount:    price.Amount,
			ProductId: price.ProductId,
			ColorId:   price.ColorId,
			SizeId:    price.SizeId,
			PriceId:   price.Id,
			Quantity:  quantityMap[key],
		})
	}
	return &TotalPriceResponse{
		TotalPrice: totalPrice,
		Items:      items,
	}
}
