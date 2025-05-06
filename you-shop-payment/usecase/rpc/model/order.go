package model

import (
	"github.com/TechwizsonORG/payment-service/entity"
	"github.com/google/uuid"
)

type OrderPayment struct {
	OrderId       uuid.UUID            `json:"orderId"`
	PaymentStatus entity.PaymentStatus `json:"payemtnStatus"`
}

func CreateOrderPayment(payment entity.Payment) *OrderPayment {
	return &OrderPayment{
		OrderId:       payment.OrderId,
		PaymentStatus: payment.Status,
	}
}

type GetOrdersPaymentRequest struct {
	OrderIds []uuid.UUID `json:"orderIds"`
}

type OrderPaymentResponse struct {
	PaymentResult []*OrderPayment `json:"paymentResult"`
}
