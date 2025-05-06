package usecase

import (
	"encoding/json"

	"github.com/TechwizsonORG/payment-service/config/model"
	"github.com/TechwizsonORG/payment-service/entity"
	"github.com/TechwizsonORG/payment-service/err"
	messagequeue "github.com/TechwizsonORG/payment-service/usecase/message_queue"
	"github.com/TechwizsonORG/payment-service/usecase/message_queue/event"
	usecaseModel "github.com/TechwizsonORG/payment-service/usecase/model"
	"github.com/TechwizsonORG/payment-service/usecase/rpc"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PaymentService struct {
	paymentRepo  PaymentRepository
	logger       zerolog.Logger
	rpcEndpoint  model.RpcServerEndpoint
	rpc          rpc.RpcInterface
	messageQueue messagequeue.MessageQueue
}

func NewPaymentService(paymentRepo PaymentRepository, logger zerolog.Logger, rpcEndpoint model.RpcServerEndpoint, rpc rpc.RpcInterface, messageQueue messagequeue.MessageQueue) *PaymentService {
	logger = logger.With().Str("usecase", "payment_service").Logger()
	return &PaymentService{
		paymentRepo:  paymentRepo,
		logger:       logger,
		rpcEndpoint:  rpcEndpoint,
		rpc:          rpc,
		messageQueue: messageQueue,
	}
}

func (p *PaymentService) CreatePayment(order usecaseModel.CreateOrder) (*entity.Payment, err.AppError) {

	orderJsonStr, parseErr := json.Marshal(order)
	if parseErr != nil {
		p.logger.Error().Err(parseErr).Msg("")
		return nil, err.NewCommonErr()
	}
	orderResJson := p.rpc.Req(p.rpcEndpoint.CreateOrder, string(orderJsonStr))
	var orderRes usecaseModel.OrderResponse
	parseErr = json.Unmarshal([]byte(orderResJson), &orderRes)
	if parseErr != nil {
		p.logger.Error().Err(parseErr).Msg("")
		return nil, err.NewCommonErr()
	}
	if orderRes.TotalPrice <= 0 {
		return nil, err.NewAppError(400, "only create payment with amount greater than 0", "only create payment with amount greater than 0", nil)
	}

	payment := entity.NewPayment(orderRes.TotalPrice, orderRes.Id, order.UserId)
	createErr := p.paymentRepo.CreatePayment(payment)
	if createErr != nil {
		p.logger.Error().Err(createErr).Msg("")
		return nil, err.NewCommonErr()
	}
	return payment, nil
}

func (s *PaymentService) ProcessPaymentCallback(paymentId uuid.UUID, transactionStatus entity.TransactionStatus) err.AppError {
	paymentStatus := entity.PaymentFailed
	if transactionStatus == entity.TransactionSuccess {
		paymentStatus = entity.PaymentSuccess
	}
	upateErr := s.paymentRepo.UpdatePaymentStatus(paymentId, paymentStatus, transactionStatus)
	if upateErr != nil {
		s.logger.Error().Err(upateErr).Msg("")
		return err.NewCommonErr()
	}
	if transactionStatus == entity.TransactionSuccess {
		s.logger.Debug().Msg("publish payment success")
		payment, getErr := s.paymentRepo.GetPaymentById(paymentId)
		if getErr != nil {
			s.logger.Error().Err(getErr).Msg("")
		} else {
			s.messageQueue.Publish(
				*messagequeue.NewDefaultExchangeConfig("you_shop", messagequeue.Topic),
				*messagequeue.NewDefaultQueueConfig("", "payment.status.success"),
				event.PaymentStatusChangedEvent{
					OrderId:   payment.OrderId,
					PaymentId: payment.Id,
					Status:    payment.Status,
				},
			)
		}
	}
	return nil
}

func (p *PaymentService) GetUserPayments(userId uuid.UUID, pageSize, pageNumber int) []*entity.Payment {
	offset := (pageNumber - 1) * pageSize
	payments, getPaymentsErr := p.paymentRepo.GetUserPayments(userId, pageSize, offset)
	if getPaymentsErr != nil {
		p.logger.Error().Err(getPaymentsErr).Msg("")
		return make([]*entity.Payment, 0)
	}
	return payments
}

func (p *PaymentService) GetPaymentsByOrderIds(orderIds []uuid.UUID) []*entity.Payment {
	payements, getErr := p.paymentRepo.GetPaymensByOrderIds(orderIds)
	if getErr != nil {
		p.logger.Error().Err(getErr).Msg("")
		return make([]*entity.Payment, 0)
	}
	return payements
}

func (p *PaymentService) GetPaymentByOrderId(orderId uuid.UUID) (*entity.Payment, err.AppError) {
	payment, getErr := p.paymentRepo.GetPaymentByOrderId(orderId)
	if getErr != nil {
		p.logger.Error().Err(getErr).Msg("")
		return nil, err.NewCommonErr()
	}
	if payment == nil {
		return nil, err.NewAppError(400, "not found payment", "not found payment", nil)
	}
	return payment, nil
}
