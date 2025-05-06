package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TechwizsonORG/order-service/api/middleware"
	"github.com/TechwizsonORG/order-service/api/model"
	orderModel "github.com/TechwizsonORG/order-service/api/model/order"
	"github.com/TechwizsonORG/order-service/api/util"
	configModel "github.com/TechwizsonORG/order-service/config/model"
	"github.com/TechwizsonORG/order-service/entity"
	"github.com/TechwizsonORG/order-service/err"
	"github.com/TechwizsonORG/order-service/usecase/order"
	usecaseOrderModel "github.com/TechwizsonORG/order-service/usecase/order/model"
	"github.com/TechwizsonORG/order-service/usecase/rpc"
	rpcModel "github.com/TechwizsonORG/order-service/usecase/rpc/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type OrderHandler struct {
	orderService order.Service
	logger       zerolog.Logger
	rpcService   rpc.RpcInterface
	rpcEndpoint  configModel.RpcServerEndpoint
}

func NewOrderHandler(orderService order.Service, logger zerolog.Logger, rpcService rpc.RpcInterface, rpcEndpoint configModel.RpcServerEndpoint) *OrderHandler {
	orderHandlerLogger := logger.With().Str("Handler", "Order").Logger()
	return &OrderHandler{
		orderService: orderService,
		logger:       orderHandlerLogger,
		rpcService:   rpcService,
		rpcEndpoint:  rpcEndpoint,
	}
}

func (o *OrderHandler) AddRoute(route *gin.RouterGroup) {
	orderRoute := route.Group("/orders")
	orderRoute.GET("/:id", middleware.AuthorizationMiddleware([]string{"admin", "guest"}, nil), o.getOrderById)
	orderRoute.GET("/user-orders", middleware.AuthorizationMiddleware([]string{"guest", "admin"}, nil), o.getCurrentUserOrders)
	orderRoute.PATCH("/:id/:status", middleware.AuthorizationMiddleware([]string{"admin"}, nil), o.updateStatus)
	orderRoute.POST("", middleware.AuthorizationMiddleware([]string{"admin", "guest"}, nil), o.createOrder)
	orderRoute.PUT("/:id", middleware.AuthorizationMiddleware([]string{"admin", "guest"}, nil), o.updateOrder)
	orderRoute.DELETE("/:id", middleware.AuthorizationMiddleware([]string{"admin"}, nil), o.deleteOrder)
}

// GetCurrentUserOrders godoc
//
//	@Tags		orders
//	@Summary	Get current user orders
//	@Param		page		query		int	false	"Page number"	default(1)
//	@Param		page_size	query		int	false	"Page size"		default(10)
//	@Success	200			{object}	model.ApiResponse{data=[]orderModel.Order}
//	@Failure	400			{object}	err.ValidationError	"Invalid request"
//	@Failure	500			{object}	err.OrderError		"Internal server error"
//	@Router		/orders/user-orders [get]
func (o *OrderHandler) getCurrentUserOrders(c *gin.Context) {
	ok, paginationErr := util.PaginationValidator(c)
	if !ok {
		c.Errors = append(c.Errors, &gin.Error{Err: paginationErr})
		return
	}
	page, pageSize := util.GetPaginationQuery(c)
	userId, _ := util.GetUserId(c)
	orders := o.orderService.GetUserOrders(userId, page, pageSize)
	result := []*orderModel.Order{}
	for _, order := range orders {
		result = append(result, orderModel.FromOrder(order))
	}
	o.getOrdersPayment(result...)
	c.JSON(200, model.SuccessResponse(result))
}

func (o *OrderHandler) getOrders(c *gin.Context) {

}

// CreateOrder godoc
//
//	@Tags		orders
//	@Summary	Create order
//	@Param		createOrder	body	usecaseOrderModel.CreateOrder	true	"Create order request body"
//	@Accept		json
//	@Success	201
//	@Failure	400	{object}	err.OrderError	"Invalid request"
//	@Failure	500	{object}	err.OrderError	"Internal server error"
//	@Router		/orders [post]
func (o *OrderHandler) createOrder(c *gin.Context) {
	var createOrder usecaseOrderModel.CreateOrder
	bindErr := c.BindJSON(&createOrder)
	if bindErr != nil {
		o.logger.Error().Err(bindErr).Msg("")
	}
	userIdStr := c.Request.Header.Get("userId")
	userId, parseErr := uuid.Parse(userIdStr)
	if parseErr != nil {
		o.logger.Error().Err(parseErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderError(401, "Couldn't parse user id", "Couldn't parse user id", nil)})
		return
	}
	createOrder.UserId = userId
	order, createOrderErr := o.orderService.CreateOrder(createOrder)
	if createOrderErr != nil {
		o.logger.Error().Err(parseErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: createOrderErr})
		return
	}
	c.Status(201)
	c.Header("Location", fmt.Sprintf("/api/v1/orders/%s", order.Id.String()))
}

// UpdateOrder godoc
//
//	@Tags		orders
//	@Summary	Update order
//	@Param		createOrder	body	orderModel.UpdateOrder	true	"Create order request body"
//	@Param		id			path	string					true	"Order ID"
//	@Accept		json
//	@Success	200	{object}	model.ApiResponse{data=orderModel.Order}
//	@Failure	400	{object}	err.OrderError	"Invalid request"
//	@Failure	500	{object}	err.OrderError	"Internal server error"
//	@Router		/orders/:id [put]
func (o *OrderHandler) updateOrder(c *gin.Context) {
	orderIdStr := c.Param("id")
	if strings.Compare("", orderIdStr) == 0 {
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderError(400, "couldn't get id", "couldn't get id", nil)})
		return
	}
	orderId, parseUuidErr := uuid.Parse(orderIdStr)
	if parseUuidErr != nil {
		o.logger.Error().Err(parseUuidErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderError(400, "couldn't parse id", "couldn't parse id", nil)})
		return
	}
	userId, getUserIdErr := util.GetUserId(c)
	if getUserIdErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderError(400, getUserIdErr.Error(), getUserIdErr.Error(), nil)})
	}
	var req orderModel.UpdateOrder
	bindErr := c.BindJSON(&req)
	if bindErr != nil {
		o.logger.Error().Err(bindErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderDefaultError(nil)})
		return
	}
	updatedOrder, updateErr := o.orderService.UpdateOrder(req.ToEntity(orderId), req.IsCancel, userId)
	if updateErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: updateErr})
		return
	}
	result := orderModel.FromOrder(*updatedOrder)
	o.getOrdersPayment(result)
	c.JSON(200, model.SuccessResponse(result))
}

func (o *OrderHandler) deleteOrder(c *gin.Context) {
	orderId, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		o.logger.Error().Err(parseErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderError(400, "couldn't parse order id", "couldn't parse order id", nil)})
		return
	}
	userId, userIdErr := util.GetUserId(c)
	if userIdErr != nil {
		o.logger.Error().Msgf("couldn't parse user id")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderDefaultError(nil)})
		return
	}
	o.orderService.DeleteOrder(orderId, userId)
}

// GetOrderById godoc
//
//	@Tags		orders
//	@Summary	Get order by Id
//	@Accept		json
//	@Success	200	{object}	model.ApiResponse{data=orderModel.Order}
//	@Failure	400	{object}	err.OrderError	"Invalid request"
//	@Failure	500	{object}	err.OrderError	"Internal server error"
//	@Router		/orders/:id [get]
func (o *OrderHandler) getOrderById(c *gin.Context) {
	orderId, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		o.logger.Error().Err(parseErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderError(400, "couldn't parse order id", "couldn't parse order id", nil)})
		return
	}
	userId, userIdErr := util.GetUserId(c)
	if userIdErr != nil {
		o.logger.Error().Msgf("couldn't parse user id")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderDefaultError(nil)})
		return
	}
	order, getOrderErr := o.orderService.GetOrder(orderId, userId)
	if getOrderErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: getOrderErr})
		return
	}

	responseOrder := orderModel.FromOrder(*order)
	o.getProducts(*order, responseOrder)
	o.getOrdersPayment(responseOrder)
	c.JSON(200, model.SuccessResponse(responseOrder))
}

func (o *OrderHandler) updateStatus(c *gin.Context) {

}

func (o *OrderHandler) getProducts(order entity.Order, orderResponse *orderModel.Order) {
	productIds := []uuid.UUID{}
	for _, item := range order.Items {
		productIds = append(productIds, item.ProductId)
	}
	req := rpcModel.GetProductByIdsRequest{
		ProductIds: productIds,
	}
	reqJson, parseErr := json.Marshal(req)
	if parseErr != nil {
		o.logger.Error().Err(parseErr).Msg("")
		return
	}
	resJson := o.rpcService.Req(o.rpcEndpoint.GetProductByIds, string(reqJson))
	var res rpcModel.GetProductByIdsResponse
	parseErr = json.Unmarshal([]byte(resJson), &res)
	if parseErr != nil {
		o.logger.Error().Err(parseErr).Msg("")
		return
	}
	productMap := map[uuid.UUID]string{}
	for _, product := range res.Products {
		productMap[product.Id] = product.Name
	}

	for _, item := range orderResponse.Items {
		if name, exist := productMap[item.ProductId]; exist {
			item.ProductName = name
		} else {
			o.logger.Warn().Msgf("Couldn't find product name for product id: %s", item.ProductId.String())
		}
	}
}

func (o *OrderHandler) getOrdersPayment(orders ...*orderModel.Order) {
	orderMap := make(map[uuid.UUID]*orderModel.Order, len(orders))
	orderIds := make([]uuid.UUID, 0, len(orderMap))
	for _, order := range orders {
		orderMap[order.Id] = order
		orderIds = append(orderIds, order.Id)
	}

	request := rpcModel.GetOrdersPaymentRequest{
		OrderIds: orderIds,
	}

	reqJson, bindErr := json.Marshal(request)
	if bindErr != nil {
		o.logger.Error().Err(bindErr).Msg("")
		return
	}

	resJson := o.rpcService.Req(o.rpcEndpoint.GetOrdersPayment, string(reqJson))
	var res rpcModel.OrderPaymentResponse
	bindErr = json.Unmarshal([]byte(resJson), &res)
	if bindErr != nil {
		o.logger.Error().Err(bindErr).Msg("")
		return
	}
	for _, paymentResult := range res.PaymentResult {
		orderMap[paymentResult.OrderId].PaymentStatus = paymentResult.PaymentStatus
	}
}
