package handler

import (
	"github.com/TechwizsonORG/order-service/api/middleware"
	"github.com/TechwizsonORG/order-service/entity"
	"github.com/TechwizsonORG/order-service/err"
	"github.com/TechwizsonORG/order-service/usecase/order"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type AdminOrderHandler struct {
	logger       zerolog.Logger
	orderService order.Service
}

func NewAdminOrderHandler(logger zerolog.Logger, orderService order.Service) *AdminOrderHandler {
	logger = logger.With().Str("Handler", "AdminOrderHandler").Logger()
	return &AdminOrderHandler{
		logger:       logger,
		orderService: orderService,
	}
}

func (a *AdminOrderHandler) AdminOrderRoute(group *gin.RouterGroup) {
	adminGroup := group.Group("/orders/admin", middleware.AuthorizationMiddleware([]string{"admin"}, nil))
	adminGroup.PATCH("/:id/confirm", a.confirmOrder)
}

// ConfirmOrder godoc
//
//	@Tags			admin-orders
//	@Summary		Confirm order
//	@Description	Confirm order by id
//	@Param			id	path	string	true	"Order ID"
//	@Router			/orders/admin/:id/confirm [patch]
//	@Success		202
//	@Failure		400	{object}	err.OrderError	"Invalid request"
//	@Failure		500	{object}	err.OrderError	"Internal server error"
func (a *AdminOrderHandler) confirmOrder(c *gin.Context) {
	orderId, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		a.logger.Error().Err(parseErr)
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewOrderError(400, "couldn't parse id", "couldn't parse id", nil)})
		return
	}
	_, updateErr := a.orderService.UpdateOrderStatus(orderId, entity.Confirmed, false)

	if updateErr != nil {
		a.logger.Error().Err(parseErr)
		c.Errors = append(c.Errors, &gin.Error{Err: updateErr})
		return
	}
	c.Status(202)
}
