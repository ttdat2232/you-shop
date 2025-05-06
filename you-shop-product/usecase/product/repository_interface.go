package product

import (
	"github.com/TechwizsonORG/product-service/entity"
	appErr "github.com/TechwizsonORG/product-service/err"
	"github.com/google/uuid"
)

type Reader interface {
	Search(query string) ([]entity.Product, appErr.ApplicationError)
	List(page int, pageSize int) ([]entity.Product, appErr.ApplicationError)
	Get(id uuid.UUID) (*entity.Product, appErr.ApplicationError)
	Count() (int, appErr.ApplicationError)
	IsSkuAlreadyExisted(sku string) (bool, appErr.ApplicationError)
	IsIdExisted(id uuid.UUID) (bool, appErr.ApplicationError)
	GetByIds(productIds []uuid.UUID) ([]entity.Product, appErr.ApplicationError)
}

type Writer interface {
	Create(product entity.Product) (entity.Product, appErr.ApplicationError)
	Update(product entity.Product) (entity.Product, appErr.ApplicationError)
	Delete(id uuid.UUID) appErr.ApplicationError
	AddProductImages([]entity.ProductImage) appErr.ApplicationError
}

type ProductRepository interface {
	Reader
	Writer
}
