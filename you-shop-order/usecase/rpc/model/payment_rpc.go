package model

import "github.com/google/uuid"

type PaymentStatus int8

const (
	PaymentPending PaymentStatus = iota + 1
	PaymentSuccess
	PaymentFailed
)

type OrderPayment struct {
	OrderId       uuid.UUID     `json:"orderId"`
	PaymentStatus PaymentStatus `json:"payemtnStatus"`
}
type GetOrdersPaymentRequest struct {
	OrderIds []uuid.UUID `json:"orderIds"`
}

type OrderPaymentResponse struct {
	PaymentResult []*OrderPayment `json:"paymentResult"`
}
