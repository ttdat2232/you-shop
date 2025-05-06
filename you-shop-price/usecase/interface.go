package usecase

import (
	"github.com/TechwizsonORG/price-service/entity"
	"github.com/TechwizsonORG/price-service/err"
	"github.com/TechwizsonORG/price-service/usecase/event"
	"github.com/TechwizsonORG/price-service/usecase/rpc/model"
	"github.com/google/uuid"
)

type Repository interface {
	Reader
	Writer
}

type Service interface {
	GetCurrentPrices(productIds []string) (map[string]float64, *err.AppError)
	CreateNewPrices(event.CreatedInventoriesEvent) ([]*entity.Price, *err.AppError)
	CreateNewPriceList(description string, currency entity.Currency) (*entity.PriceList, *err.AppError)
	UpdatePrice(productId, colorId, sizeId uuid.UUID, price float64) (bool, *err.AppError)
	GetTotalPrice(model.TotalPriceRequest) (float64, []entity.Price, *err.AppError)
}

type Reader interface {
	GetCurrentPricesByProductIds([]uuid.UUID) ([]entity.Price, error)
	GetCurrentPrices([]*model.OrderItem) ([]entity.Price, error)
	GetDefaultPriceList() *entity.PriceList
	GetPrice(productId uuid.UUID, colorId uuid.UUID, sizeId uuid.UUID) (*entity.Price, error)
}

type Writer interface {
	AddNewPrice(entity.Price) (*entity.Price, error)
	AddNewPrices([]*entity.Price) error
	AddNewPriceList(entity.PriceList) (*entity.PriceList, error)
	UpdatePrice(productId, sizeId, colorId uuid.UUID, price float64) (bool, error)
}
