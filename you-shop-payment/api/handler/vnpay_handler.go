package handler

import (
	"github.com/TechwizsonORG/payment-service/api/middleware"
	"github.com/TechwizsonORG/payment-service/api/model"
	vnpayModel "github.com/TechwizsonORG/payment-service/api/model/payment"
	"github.com/TechwizsonORG/payment-service/api/util"
	"github.com/TechwizsonORG/payment-service/err"
	"github.com/TechwizsonORG/payment-service/infrastructure"
	"github.com/TechwizsonORG/payment-service/usecase"
	usecaseModel "github.com/TechwizsonORG/payment-service/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type VnpayHandler struct {
	vnpayService   *infrastructure.VnpayService
	paymentService usecase.Service
	logger         zerolog.Logger
}

func NewVnpayHandler(vnpayService *infrastructure.VnpayService, paymentService usecase.Service, logger zerolog.Logger) *VnpayHandler {
	logger = logger.With().Str("handler", "vnpay").Logger()
	return &VnpayHandler{
		vnpayService:   vnpayService,
		paymentService: paymentService,
		logger:         logger,
	}
}

func (v *VnpayHandler) AddVnpayRoute(c *gin.RouterGroup) {
	vnpayGroup := c.Group("/payments/vnpay")
	vnpayGroup.POST("/pay", middleware.AuthorizationMiddleware([]string{"admin", "guest"}, nil), v.getVnpayPaymentUrl)
	vnpayGroup.GET("/pay/:orderId", middleware.AuthorizationMiddleware([]string{"admin", "guest"}, nil), v.reGetVnpayPaymentUrl)
	vnpayGroup.GET("/ipn", v.vnpayCallback)
}

// VnPayCallBack godoc
//
//	@Summary	Handle VnPay Callback IPN
//	@Tags		vnpay
//	@Router		/payments/vnpay/ipn [get]
func (v *VnpayHandler) vnpayCallback(c *gin.Context) {
	v.logger.Info().Msg("vnpay callback")
	paymentId, transactionStatus, callBackerr := v.vnpayService.VnpPayHandleCallback(c.Request)
	res := map[string]any{}
	res["RspCode"] = "99"
	res["Message"] = "Falied"
	if callBackerr != nil {
		c.JSON(200, res)
		return
	}
	processErr := v.paymentService.ProcessPaymentCallback(*paymentId, transactionStatus)
	if processErr != nil {
		c.JSON(200, res)
		return
	}
	res["RspCode"] = "00"
	res["Message"] = "Success"
	c.JSON(200, res)
}

// GetVnPayPaymentUrl godoc
//
//	@Summary		Generate VnPay Payment URL
//	@Tags			vnpay
//	@Router			/payments/vnpay/pay [post]
//	@Description	Generate VnPay Payment URL and also creating order
//	@Accept			json
//	@Param			request	body	usecaseModel.CreateOrder	true	"Create Order Request"
//	@Produce		json
//	@Success		200	{object}	model.ApiResponse{data=vnpayModel.PaymentUrlResponse}
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
func (v *VnpayHandler) getVnpayPaymentUrl(c *gin.Context) {
	var createOrderReq usecaseModel.CreateOrder
	bindErr := c.BindJSON(&createOrderReq)
	if bindErr != nil {
		v.logger.Error().Err(bindErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewAppError(400, "couldn't parse request JSON", "couldn't parse request JSON", nil)})
		return
	}

	userId, _ := util.GetUserId(c)
	createOrderReq.UserId = userId
	payment, createErr := v.paymentService.CreatePayment(createOrderReq)
	if createErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: createErr})
		return
	}
	url, getErr := v.vnpayService.CreatePaymentUrl(*payment, c.Request.RemoteAddr)
	if getErr != nil {
		v.logger.Error().Err(getErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewAppError(400, "failed when generating url", "failed when generating url", nil)})
		return
	}
	c.JSON(200, model.SuccessResponse(vnpayModel.PaymentUrlResponse{
		Url: url,
	}))
}

// ReGetVnpayPaymentUrl godoc
//
//	@Summary		Regenerate VnPay Payment URL
//	@Description	Regenerate VnPay Payment URL for existed payment via order Id
//	@Tags			vnpay
//	@Router			/payments/vnpay/pay/:orderId [get]
//	@Param			orderId	query	string	false	"Order Id"
//	@Produce		json
//	@Success		200	{object}	model.ApiResponse{data=vnpayModel.PaymentUrlResponse}
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
func (v *VnpayHandler) reGetVnpayPaymentUrl(c *gin.Context) {
	orderId, parseErr := uuid.Parse(c.Param("orderId"))

	if parseErr != nil {
		v.logger.Error().Err(parseErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewAppError(400, "failed when parse order Id", "failed when parse order Id", nil)})
		return
	}

	url, getErr := v.vnpayService.CreatePaymentUrlByOrderId(orderId, c.Request.RemoteAddr)

	if getErr != nil {
		v.logger.Error().Err(getErr).Msg("")
		return
	}
	c.JSON(200, model.SuccessResponse(vnpayModel.PaymentUrlResponse{
		Url: url,
	}))
}
