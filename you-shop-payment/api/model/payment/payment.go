package payment

import (
	"github.com/TechwizsonORG/payment-service/entity"
	"github.com/google/uuid"
)

type PaymentResponse struct {
	Id      uuid.UUID            `json:"id"`
	OrderId uuid.UUID            `json:"orderId"`
	UserId  uuid.UUID            `json:"userId"`
	Status  entity.PaymentStatus `json:"status"`
	Amount  float64              `json:"amount"`
}

func FromPayementEntity(payment *entity.Payment) *PaymentResponse {
	return &PaymentResponse{
		Id:      payment.Id,
		OrderId: payment.OrderId,
		UserId:  payment.UserId,
		Status:  payment.Status,
		Amount:  payment.Amount,
	}
}
