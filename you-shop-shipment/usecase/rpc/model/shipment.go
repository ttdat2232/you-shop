package model

import "github.com/google/uuid"

type CreateShipment struct {
	OrderId      uuid.UUID `json:"orderId"`
	PhoneNumber  string    `json:"phoneNumber"`
	ReceiverName string    `json:"receiverName"`
	ProvinceId   int       `json:"provinceId"`
	DistrictId   int       `json:"districtId"`
	WardId       int       `json:"wardId"`
}

type CreateShipmentResponse struct {
	Id           uuid.UUID `json:"id"`
	Address      string    `json:"address"`
	DevliveryFee float64   `json:"devliveryFee"`
}
