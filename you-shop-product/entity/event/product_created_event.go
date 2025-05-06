package event

import (
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/google/uuid"
)

type ProductCreatedEvent struct {
	ProductId uuid.UUID `json:"productId"`
	Price     float64   `json:"price"`
}

func NewProductCreatedEvent(product entity.Product, sizeId, colorId uuid.UUID, price float64) *ProductCreatedEvent {
	return &ProductCreatedEvent{
		ProductId: product.Id,
		Price:     price,
	}
}
