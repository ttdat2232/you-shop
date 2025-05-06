package model

import (
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/google/uuid"
)

type GetProductByIdsResponse struct {
	Products []Product `json:"products"`
}

type Product struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func From(productEntities []entity.Product) *GetProductByIdsResponse {
	products := []Product{}
	for _, product := range productEntities {
		products = append(products, Product{
			Name: product.Name,
			Id:   product.Id,
		})
	}
	return &GetProductByIdsResponse{
		Products: products,
	}
}
