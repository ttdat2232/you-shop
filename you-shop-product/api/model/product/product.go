package product

import (
	"time"

	"github.com/TechwizsonORG/product-service/entity"
	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Weight      float32              `json:"weight"`
	Weight_unit string               `json:"weight_unit"`
	Quantity    int                  `json:"quantity"`
	Price       float64              `json:"price"`
	Images      []string             `json:"images"`
	Thumbnail   string               `json:"thumbnail"`
	Status      entity.ProductStatus `json:"status"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func FromEntity(product entity.Product) Product {
	return Product{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Status:      product.Status,
		Thumbnail:   product.Thumbnail,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
