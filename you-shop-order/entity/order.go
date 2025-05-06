package entity

import (
	"github.com/TechwizsonORG/order-service/util"
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

func (ostatus OrderStatus) String() string {
	switch ostatus {
	case Pending:
		return "Pending"
	case Confirmed:
		return "Confirmed"
	case Processing:
		return "Processing"
	case Shipped:
		return "Shipped"
	case OutForDelivery:
		return "Out For Delivery"
	case Delivered:
		return "Delivered"
	case Refunded:
		return "Refunded"
	case Returned:
		return "Returned"
	case Failed:
		return "Failed"
	case Canceled:
		return "Canceled"
	case Completed:
		return "Completed"
	default:
		return "Unknown"
	}
}

type Order struct {
	AuditEntity
	Description string
	OrderCode   string
	TotalPrice  float64
	Status      OrderStatus
	OwnerId     uuid.UUID
	Items       []*OrderItem
}

// Create order with status Pending
func NewOrder(description, orderCode string, totalPrice float64, ownerId uuid.UUID, items []*OrderItem) *Order {
	orderId := uuid.New()
	for _, value := range items {
		value.OrderId = orderId
	}
	return &Order{
		Description: description,
		TotalPrice:  totalPrice,
		Status:      Pending,
		OwnerId:     ownerId,
		OrderCode:   orderCode,
		AuditEntity: AuditEntity{
			CreatedAt: util.GetCurrentUtcTime(7),
			UpdatedAt: util.GetCurrentUtcTime(7),
			IsDeleted: false,
			Id:        orderId,
		},
		Items: items,
	}
}
