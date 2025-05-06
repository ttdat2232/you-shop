package handler

import (
	"github.com/TechwizsonORG/payment-service/api/middleware"
	"github.com/TechwizsonORG/payment-service/api/model"
	"github.com/TechwizsonORG/payment-service/api/model/payment"
	"github.com/TechwizsonORG/payment-service/api/util"
	"github.com/TechwizsonORG/payment-service/err"
	"github.com/TechwizsonORG/payment-service/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type PaymentHandler struct {
	paymentService usecase.Service
	logger         zerolog.Logger
}

func NewPaymentHandler(paymentService usecase.Service, logger zerolog.Logger) *PaymentHandler {
	logger = logger.With().Str("handler", "payment_handler").Logger()
	return &PaymentHandler{
		paymentService: paymentService,
		logger:         logger,
	}
}

func (p *PaymentHandler) AddPaymentRoute(c *gin.RouterGroup) {
	route := c.Group("/payments")
	route.GET("/user-payments", middleware.AuthorizationMiddleware([]string{"guest", "admin"}, nil), p.getUserPayments)
}

// GetUserPayments godoc
//
//	@Summary	Get current user payments
//	@Success	200			{object}	model.ApiResponse{data=[]payment.PaymentResponse}
//	@Param		page		query		int	false	"Page number"	default(1)
//	@Param		page_size	query		int	false	"Page size"		default(10)
//	@Tags		payments
//	@Router		/payments/user-payments [get]
func (p *PaymentHandler) getUserPayments(c *gin.Context) {
	p.logger.Debug().Msg("get user payments")
	userId, getUserIdErr := util.GetUserId(c)
	pageSize, pageNumber := util.GetPaginationQuery(c)

	if getUserIdErr != nil {
		p.logger.Error().Err(getUserIdErr).Msg("failed to get user id")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewAppError(401, "unauthorized", "unauthorized", nil)})
		return
	}

	payments := p.paymentService.GetUserPayments(userId, pageSize, pageNumber)
	res := make([]*payment.PaymentResponse, len(payments))
	for i, paymentEntity := range payments {
		res[i] = payment.FromPayementEntity(paymentEntity)
	}
	c.JSON(200, model.SuccessResponse(res))
}
