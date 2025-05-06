package order

import (
	"time"

	"github.com/TechwizsonORG/order-service/entity"
	"github.com/TechwizsonORG/order-service/usecase/rpc/model"
	"github.com/google/uuid"
)

type Order struct {
	Id            uuid.UUID           `json:"id"`
	Description   string              `json:"description"`
	TotalPrice    float64             `json:"totalPrice"`
	CreatedAt     time.Time           `json:"createdAt"`
	Status        entity.OrderStatus  `json:"status"`
	Items         []*Item             `json:"items,omitempty"`
	PaymentStatus model.PaymentStatus `json:"paymentStatus"`
}

type Item struct {
	ProductId   uuid.UUID `json:"productId"`
	ProductName string    `json:"productName"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
}

func FromOrderItem(orderItems []*entity.OrderItem) []*Item {
	items := []*Item{}
	for _, item := range orderItems {
		items = append(items, &Item{
			ProductId: item.ProductId,
			Price:     item.Price,
			Quantity:  item.Quantity,
		})
	}
	return items
}

func FromOrder(o entity.Order) *Order {
	orderItems := FromOrderItem(o.Items)
	return &Order{
		Id:          o.Id,
		Description: o.Description,
		TotalPrice:  o.TotalPrice,
		CreatedAt:   o.CreatedAt,
		Status:      o.Status,
		Items:       orderItems,
	}
}
