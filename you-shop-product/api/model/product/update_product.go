package product

import (
	"github.com/TechwizsonORG/product-service/entity"
)

type UpdateProduct struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Status      entity.ProductStatus `json:"status"`
	Sku         string               `json:"sku"`
	UserManual  string               `json:"userManual"`
}
