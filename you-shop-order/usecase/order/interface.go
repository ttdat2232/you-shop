package order

import (
	"github.com/TechwizsonORG/order-service/entity"
	"github.com/TechwizsonORG/order-service/err"
	"github.com/TechwizsonORG/order-service/usecase/order/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetOrders(pageIndex int, pageSize int) ([]entity.Order, error)
	GetOrderById(uuid.UUID) (*entity.Order, error)
	GetOrdersByUserId(userId uuid.UUID, pageIndex int, pageSize int) ([]entity.Order, error)
	IsOrderExistById(uuid.UUID) (bool, error)
	CreateOrder(*entity.Order) error
	UpdateOrder(*entity.Order) error
	DeleteOrder(*entity.Order) error
}
type Service interface {
	CreateOrder(createOrder model.CreateOrder) (*entity.Order, err.ApplicationError)
	UpdateOrder(updateOrder entity.Order, isCancel bool, ownerId uuid.UUID) (*entity.Order, err.ApplicationError)
	UpdateOrderStatus(orderId uuid.UUID, status entity.OrderStatus, isUpdatedByDeliver bool) (*entity.Order, err.ApplicationError)
	GetOrder(orderId, ownerId uuid.UUID) (*entity.Order, err.ApplicationError)
	GetOrders() []entity.Order
	GetUserOrders(userId uuid.UUID, page, pageSize int) []entity.Order
	DeleteOrder(orderId, ownerId uuid.UUID) err.ApplicationError
}
