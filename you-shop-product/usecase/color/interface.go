package color

import (
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/TechwizsonORG/product-service/err"
	"github.com/google/uuid"
)

type ColorUsecase interface {
	AddColor(name string) (*entity.Color, err.ApplicationError)
	GetColors() []entity.Color
	GetById(uuid.UUID) (entity.Color, err.ApplicationError)
}

type Repository interface {
	AddColor(*entity.Color) error
	GetColors() ([]entity.Color, error)
	GetBydId(uuid.UUID) (*entity.Color, error)
}
