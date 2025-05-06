package model

import (
	"github.com/TechwizsonORG/order-service/usecase/order/model"
	"github.com/google/uuid"
)

type TotalPriceRequest struct {
	Items []*OrderItem `json:"items"`
}

type OrderItem struct {
	ProductId uuid.UUID `json:"productId"`
	SizeId    uuid.UUID `json:"sizeId"`
	ColorId   uuid.UUID `json:"colorId"`
	Quantity  int       `json:"quantity"`
}

func FromOrderItems(items []*model.CreateOrderItem) []*OrderItem {
	result := make([]*OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, &OrderItem{
			ProductId: item.ProductId,
			ColorId:   item.ColorId,
			SizeId:    item.SizeId,
			Quantity:  item.Quantity,
		})
	}
	return result
}
