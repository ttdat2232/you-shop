package usecase

import (
	"github.com/TechwizsonORG/shipment-service/entity"
	"github.com/TechwizsonORG/shipment-service/err"
	"github.com/google/uuid"
)

type ShipmentUsecase interface {
	CreateShipment(entity.Shipment) (entity.Shipment, err.AppError)
	UpdateStatus(vendorId, shipmentId uuid.UUID) err.AppError
}
