package order

import (
	"encoding/json"
	"fmt"

	configModel "github.com/TechwizsonORG/order-service/config/model"
	"github.com/TechwizsonORG/order-service/entity"
	"github.com/TechwizsonORG/order-service/err"
	"github.com/TechwizsonORG/order-service/usecase/messagequeue"
	"github.com/TechwizsonORG/order-service/usecase/messagequeue/event"
	"github.com/TechwizsonORG/order-service/usecase/order/model"
	"github.com/TechwizsonORG/order-service/usecase/rpc"
	rpcModel "github.com/TechwizsonORG/order-service/usecase/rpc/model"
	"github.com/TechwizsonORG/order-service/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type OrderService struct {
	repo        Repository
	rpcService  rpc.RpcInterface
	rpcEndpoint configModel.RpcServerEndpoint
	logger      zerolog.Logger
	msq         messagequeue.MessageQueue
}

func NewOrderService(rpcService rpc.RpcInterface, rpcEndpoint configModel.RpcServerEndpoint, repo Repository, msq messagequeue.MessageQueue, logger zerolog.Logger) *OrderService {
	serviceLogger := logger.With().Str("Order", "Service").Logger()
	return &OrderService{
		rpcService:  rpcService,
		rpcEndpoint: rpcEndpoint,
		logger:      serviceLogger,
		repo:        repo,
		msq:         msq,
	}
}

func (o *OrderService) CreateOrder(createOrder model.CreateOrder) (*entity.Order, err.ApplicationError) {

	if quantityCheckErr := o.checkProductQuantity(createOrder.Items); quantityCheckErr != nil {
		return nil, quantityCheckErr
	}
	totalPrice, items, getTotalPriceErr := o.getTotalPrice(createOrder.Items)
	if getTotalPriceErr != nil {
		return nil, getTotalPriceErr
	}

	if totalPrice == 0 {
		return nil, err.NewOrderDefaultError(nil)
	}

	order := entity.NewOrder(createOrder.Description, generateOrderCode(), totalPrice, createOrder.UserId, items)
	orderJson, _ := json.Marshal(order)
	o.logger.Debug().Msgf("new order: %s", orderJson)
	createError := o.repo.CreateOrder(order)
	if createError != nil {
		o.logger.Error().Err(createError).Msg("")
		return nil, err.NewOrderDefaultError(nil)
	}
	return order, nil
}
func generateOrderCode() string {
	util.RandString(10)
	return fmt.Sprintf("%s-%s%d", "ORDER", util.RandString(10), util.GetCurrentUtcTime(7).UnixNano())
}

func (o *OrderService) UpdateOrder(updateOrder entity.Order, isCancel bool, ownerId uuid.UUID) (*entity.Order, err.ApplicationError) {
	ok, isExistErr := o.repo.IsOrderExistById(updateOrder.Id)
	if isExistErr != nil {
		o.logger.Error().Err(isExistErr).Msg("")
		return nil, err.NewOrderDefaultError(nil)
	}

	if !ok {
		return nil, err.NewValidationError("couldn't find order", "couldn't find order", []err.ValidationErrorField{{Field: "id", Message: updateOrder.Id.String()}})
	}

	order, getOrderErr := o.repo.GetOrderById(updateOrder.Id)
	if getOrderErr != nil {
		o.logger.Error().Err(getOrderErr).Msg("")
		return nil, err.NewOrderDefaultError(nil)
	}
	if ownerId != order.OwnerId {
		return nil, err.NewValidationError("couldn't modify other's order", "couldn't modify other's order", nil)
	}
	if order.Status != entity.Pending && order.Status != entity.Processing {
		return nil, err.NewValidationError("couldn't update order with current status", "couldn't update order with current status", []err.ValidationErrorField{{"status", order.Status.String()}})
	}
	if isCancel {
		order.Status = entity.Canceled
	}

	updateErr := o.repo.UpdateOrder(order)
	if updateErr != nil {
		o.logger.Error().Err(updateErr).Msg("")
		return nil, err.NewOrderDefaultError(nil)
	}
	return order, nil
}

func (o *OrderService) UpdateOrderStatus(orderId uuid.UUID, status entity.OrderStatus, isUpdatedByDeliver bool) (*entity.Order, err.ApplicationError) {
	currentOrder, getErr := o.repo.GetOrderById(orderId)
	if getErr != nil {
		o.logger.Error().Err(getErr).Msg("")
		return nil, nil
	}
	if currentOrder == nil {
		return nil, err.NewOrderError(404, "couldn't found order", "couldn't found order", nil)
	}

	if ok, msg := o.isValidStatusTransition(currentOrder.Status, status, isUpdatedByDeliver); !ok {
		return nil, err.NewOrderError(400, msg, msg, nil)
	}
	currentOrder.Status = status
	updateErr := o.repo.UpdateOrder(currentOrder)
	if updateErr != nil {
		o.logger.Error().Err(updateErr).Msg("")
		return nil, err.NewOrderDefaultError(nil)
	}
	go o.msq.Publish(
		*messagequeue.NewDefaultExchangeConfig("you_shop", messagequeue.Topic),
		*messagequeue.NewDefaultQueueConfig("", "order.updated"),
		*event.FromOrderEntity(currentOrder),
	)
	return currentOrder, nil
}

func (o *OrderService) isValidStatusTransition(currentStatus, newStatus entity.OrderStatus, isUpdatedByDeilver bool) (bool, string) {
	msg := "status trannsition failed"
	ok := false

	if currentStatus == newStatus {
		return true, ""
	}
	switch currentStatus {
	case entity.Pending:
		ok = newStatus == entity.Confirmed || newStatus == entity.Canceled || newStatus == entity.Failed
	case entity.Confirmed:
		ok = newStatus == entity.Processing || newStatus == entity.Canceled
	case entity.Processing:
		ok = newStatus == entity.Shipped || newStatus == entity.Canceled
	case entity.Shipped:
		ok = isUpdatedByDeilver && (newStatus == entity.OutForDelivery || newStatus == entity.Returned)
	case entity.OutForDelivery:
		ok = isUpdatedByDeilver && (newStatus == entity.Delivered || newStatus == entity.Returned)
	case entity.Delivered:
		ok = newStatus == entity.Refunded || newStatus == entity.Completed
	case entity.Refunded, entity.Returned, entity.Failed, entity.Canceled, entity.Completed:
		ok = false
	default:
		ok = false
	}

	if ok {
		return ok, ""
	}
	return ok, msg
}

func (o *OrderService) GetOrder(orderId, ownerId uuid.UUID) (*entity.Order, err.ApplicationError) {
	ok, checkExistErr := o.repo.IsOrderExistById(orderId)
	if checkExistErr != nil {
		o.logger.Error().Err(checkExistErr).Msg("")
		return nil, err.NewOrderDefaultError(nil)
	}

	if !ok {
		return nil, err.NewValidationError("couldn't found order", "couldn't found order", nil)
	}

	order, getOrderErr := o.repo.GetOrderById(orderId)
	if getOrderErr != nil {
		o.logger.Error().Err(getOrderErr).Msg("")
		return nil, err.NewOrderDefaultError(nil)
	}
	if order.OwnerId != ownerId {
		return nil, err.NewValidationError("couldn't found order", "couldn't found order", nil)
	}

	return order, nil
}
func (o *OrderService) GetOrders() []entity.Order {
	return nil
}
func (o *OrderService) GetUserOrders(userId uuid.UUID, page, pageSize int) []entity.Order {
	defaultResult := []entity.Order{}
	pageIndex := page - 1
	orders, getErr := o.repo.GetOrdersByUserId(userId, pageIndex, pageSize)
	if getErr != nil {
		o.logger.Error().Err(getErr).Msg("")
		return defaultResult
	}
	return orders
}

func (o *OrderService) getTotalPrice(orderItems []*model.CreateOrderItem) (float64, []*entity.OrderItem, err.ApplicationError) {
	jsonReq, parseJsonErr := json.Marshal(&rpcModel.TotalPriceRequest{
		Items: rpcModel.FromOrderItems(orderItems),
	},
	)
	if parseJsonErr != nil {
		o.logger.Error().Err(parseJsonErr).Msg("")
		return 0, nil, err.NewOrderDefaultError(nil)
	}
	jsonRes := o.rpcService.Req(o.rpcEndpoint.GetTotalPrice, string(jsonReq))
	var res rpcModel.TotalPriceResponse
	parseJsonErr = json.Unmarshal([]byte(jsonRes), &res)
	if parseJsonErr != nil {
		o.logger.Error().Err(parseJsonErr).Msg("")
		return 0, nil, err.NewOrderDefaultError(nil)
	}
	orderItemEntities := []*entity.OrderItem{}
	for _, item := range res.Items {
		orderItemEntities = append(orderItemEntities, &entity.OrderItem{
			Quantity:  item.Quantity,
			ProductId: item.ProductId,
			SizeId:    item.SizeId,
			ColorId:   item.ColorId,
			Price:     item.Amount,
			PriceId:   item.PriceId,
		})
	}
	return res.TotalPrice, orderItemEntities, nil
}

func (o *OrderService) checkProductQuantity(orderItems []*model.CreateOrderItem) err.ApplicationError {
	for _, item := range orderItems {
		req := &rpcModel.CheckProductQuantity{
			ProductId:       item.ProductId,
			ColorId:         item.ColorId,
			SizeId:          item.SizeId,
			RequireQuantity: item.Quantity,
		}
		jsonReq, parseJsonErr := json.Marshal(req)
		if parseJsonErr != nil {
			o.logger.Error().Err(parseJsonErr).Msg("")
			return err.NewOrderDefaultError(nil)
		}
		jsonRes := o.rpcService.Req(o.rpcEndpoint.CheckProductQuantity, string(jsonReq))
		var res rpcModel.CheckProductQuantityResponse
		parseJsonErr = json.Unmarshal([]byte(jsonRes), &res)
		if parseJsonErr != nil {
			o.logger.Error().Err(parseJsonErr).Msg("")
			return err.NewOrderDefaultError(nil)
		}
		if !res.IsEnough {
			title := "Not enough quantity"
			detail := fmt.Sprintf("Not enough quantity for product: %s", item.ProductId)
			return err.NewOrderError(400, title, detail, nil)
		}
	}
	return nil
}

func (o *OrderService) DeleteOrder(orderId, ownerId uuid.UUID) err.ApplicationError {
	ok, checkExistErr := o.repo.IsOrderExistById(orderId)
	if checkExistErr != nil {
		o.logger.Error().Err(checkExistErr).Msg("")
		return err.NewOrderDefaultError(nil)
	}

	if !ok {
		return err.NewValidationError("couldn't found order", "couldn't found order", nil)
	}

	order, getOrderErr := o.repo.GetOrderById(orderId)
	if getOrderErr != nil {
		o.logger.Error().Err(getOrderErr).Msg("")
		return err.NewOrderDefaultError(nil)
	}

	if order.OwnerId != ownerId {
		return err.NewValidationError("couldn't found order", "couldn't found order", nil)
	}

	order.IsDeleted = true
	order.DeletedAt = util.GetCurrentUtcTime(7)
	order.UpdatedAt = util.GetCurrentUtcTime(7)
	updateErr := o.repo.UpdateOrder(order)
	if updateErr != nil {
		o.logger.Error().Err(updateErr).Msg("")
		return err.NewOrderDefaultError(nil)
	}
	return nil
}
