package size

import (
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/TechwizsonORG/product-service/err"
)

type SizeUseCase interface {
	AddSize(name string) (*entity.Size, err.ApplicationError)
	GetSizes() []entity.Size
}

type SizeRepository interface {
	AddSize(*entity.Size) error
	GetSizes() ([]entity.Size, error)
}
