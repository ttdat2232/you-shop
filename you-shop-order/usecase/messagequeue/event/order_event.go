package event

import (
	"github.com/TechwizsonORG/order-service/entity"
	"github.com/google/uuid"
)

type OrderStatus int8

const (
	Pending        OrderStatus = iota + 1 // Created but not yet processed
	Confirmed                             // Confirmed by the admin
	Processing                            // Being prepared or processed
	Shipped                               // The order has been dispatched to the courier
	OutForDelivery                        // The order is on its way to be delivered
	Delivered                             // The order has been delivered
	Refunded                              // The order has been refunded to the customer
	Returned                              // The order has been returned by the customer
	Failed                                // Payment failed or any other failure
	Canceled                              // The order canceled by the user or system
	Completed                             // The order has been fully processed, delivered, and closed
)

type UpdatedOrderEvent struct {
	Id          uuid.UUID           `json:"id"`
	Description string              `json:"description"`
	OrderCode   string              `json:"orderCode"`
	TotalPrice  float64             `json:"totalPrice"`
	Status      OrderStatus         `json:"status"`
	OwnerId     uuid.UUID           `json:"ownerId"`
	Items       []*UpdatedOrderItem `json:"items`
}

type UpdatedOrderItem struct {
	Quantity  int       `json:"quantity"`
	ProductId uuid.UUID `json:"productId"`
	ColorId   uuid.UUID `json:"colorId"`
	SizeId    uuid.UUID `json:"sizeId"`
	OrderId   uuid.UUID `json:"orderId"`
	Price     float64   `json:"price"`
}

func FromOrderEntity(order *entity.Order) *UpdatedOrderEvent {
	items := make([]*UpdatedOrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &UpdatedOrderItem{
			Quantity:  item.Quantity,
			ProductId: item.ProductId,
			ColorId:   item.ColorId,
			SizeId:    item.SizeId,
			OrderId:   item.OrderId,
			Price:     item.Price,
		})
	}
	return &UpdatedOrderEvent{
		Id:          order.Id,
		Description: order.Description,
		OrderCode:   order.OrderCode,
		TotalPrice:  order.TotalPrice,
		Status:      OrderStatus(order.Status),
		OwnerId:     order.OwnerId,
		Items:       items,
	}
}
