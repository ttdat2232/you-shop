package product

import (
	"mime/multipart"

	"github.com/TechwizsonORG/product-service/entity"
	appErr "github.com/TechwizsonORG/product-service/err"
	"github.com/google/uuid"
)

type UseCase interface {
	SearchProducts(query string) []entity.Product
	GetProducts(page int, pageSize int) (count int, products []entity.Product)
	GetProductByIds([]uuid.UUID) []entity.Product
	GetProduct(id string) (product *entity.Product, appErr appErr.ApplicationError)
	CreateProduct(name, description, sku, userManual string, productImages map[string]*multipart.File, thumbnailImage *multipart.FileHeader) (product *entity.Product, appErr appErr.ApplicationError)
	UpdateProduct(id uuid.UUID, name, description, sku string, status entity.ProductStatus, userManual string) (product *entity.Product, appErr appErr.ApplicationError)
	DeleteProduct(id uuid.UUID) appErr.ApplicationError
	CheckProductQuantity(productId, colorId, sizeId uuid.UUID, requireQuantity int) (bool, appErr.ApplicationError)
	GetQuantity(productId, sizeId, colorId uuid.UUID) (quantity int, appErr appErr.ApplicationError)
	UploadProductColor(productId, colorId uuid.UUID, file *multipart.FileHeader) appErr.ApplicationError
}
