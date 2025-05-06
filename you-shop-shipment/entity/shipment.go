package entity

import (
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus int8

const ()

type ShipmentItem struct {
	Name   string
	Height float64
	Weight float64
	Width  float64
}
type Shipment struct {
	AuditEntity
	Carrier          string
	OrderId          uuid.UUID
	Status           ShipmentStatus
	VendorId         uuid.UUID
	DeliveryTime     time.Time
	DeliveryTotalFee float64
	AddressLine      string
	WardId           int64
	ProvinceId       int64
	Items            []*ShipmentItem
}
