package job

import (
	"encoding/json"

	apiModel "github.com/TechwizsonORG/order-service/api/model/order"
	"github.com/TechwizsonORG/order-service/background"
	"github.com/TechwizsonORG/order-service/entity"
	"github.com/TechwizsonORG/order-service/usecase/messagequeue"
	"github.com/TechwizsonORG/order-service/usecase/messagequeue/event"
	"github.com/TechwizsonORG/order-service/usecase/order"
	"github.com/TechwizsonORG/order-service/usecase/order/model"
	"github.com/TechwizsonORG/order-service/usecase/rpc"
	"github.com/rs/zerolog"
)

type Job struct {
	logger zerolog.Logger
}

func NewJob(logger zerolog.Logger) *Job {
	jobLogger := logger.
		With().
		Str("Background", "Job").
		Logger()
	return &Job{logger: jobLogger}
}

func (j *Job) CreateOrder(rpcService rpc.RpcInterface, orderService order.Service) background.JobFunc {
	return func() {
		j.logger.Debug().Msg("Create order from background job")
		rpcService.NewRpcQueue("create_order", func(data string) string {
			var createOrderModel model.CreateOrder
			bindErr := json.Unmarshal([]byte(data), &createOrderModel)
			if bindErr != nil {
				j.logger.Error().Err(bindErr).Msg("")
				return ""
			}
			order, createErr := orderService.CreateOrder(createOrderModel)
			if createErr != nil {
				j.logger.Error().Err(createErr).Msg("")
				return ""
			}
			orderRes := apiModel.FromOrder(*order)
			orderResJson, bindJsonErr := json.Marshal(orderRes)
			if bindJsonErr != nil {
				return ""
			}
			return string(orderResJson)
		})
	}
}

func (j *Job) HandlePaymentStatusChangedEvent(msq messagequeue.MessageQueue, orderService order.Service) background.JobFunc {
	return func() {
		j.logger.Debug().Msg("Handle payment status changed event from background job")
		msq.Consume(
			*messagequeue.NewDefaultExchangeConfig("you_shop", messagequeue.Topic),
			*messagequeue.NewDefaultQueueConfig("payment_status_changed_event", "payment.status.*"),
			func(data string) error {
				j.logger.Debug().Msgf("Received data: %s", data)
				var paymentStatusChangedEvent event.PaymentStatusChangedEvent
				bindErr := json.Unmarshal([]byte(data), &paymentStatusChangedEvent)
				if bindErr != nil {
					j.logger.Error().Err(bindErr).Msg("")
					return bindErr
				}

				if paymentStatusChangedEvent.Status != event.PaymentSuccess {
					j.logger.Debug().Msg("Payment status is not success, skipping order update")
					return nil
				}
				j.logger.Debug().Msgf("Payment status is success, updating order with ID: %s", paymentStatusChangedEvent.OrderId.String())
				_, updateErr := orderService.UpdateOrderStatus(paymentStatusChangedEvent.OrderId, entity.Confirmed, false)
				if updateErr != nil {
					j.logger.Error().Err(updateErr).Msg("")
					return updateErr
				}
				return nil
			},
		)
	}
}
