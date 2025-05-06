package entity

import (
	"time"

	"github.com/TechwizsonORG/price-service/util"
	"github.com/google/uuid"
)

type Price struct {
	AuditEntity
	Amount      float64
	ValidFrom   time.Time
	ValidTo     time.Time
	MinQuantity int
	ProductId   uuid.UUID
	ColorId     uuid.UUID
	SizeId      uuid.UUID
	IsActive    bool
	PriceListId uuid.UUID
}

func NewPrice(amount float64, productId, colorId, sizeId uuid.UUID) *Price {
	return &Price{
		AuditEntity: AuditEntity{
			Id:        uuid.New(),
			CreatedAt: util.GetCurrentUtcTime(7),
			UpdatedAt: util.GetCurrentUtcTime(7),
		},
		Amount:    amount,
		ProductId: productId,
		ColorId:   colorId,
		SizeId:    sizeId,
		IsActive:  true,
		ValidFrom: util.GetCurrentUtcTime(7),
	}
}
