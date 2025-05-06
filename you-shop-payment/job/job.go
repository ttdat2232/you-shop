package job

import (
	"encoding/json"

	"github.com/TechwizsonORG/payment-service/background"
	"github.com/TechwizsonORG/payment-service/usecase"
	"github.com/TechwizsonORG/payment-service/usecase/rpc"
	"github.com/TechwizsonORG/payment-service/usecase/rpc/model"
	"github.com/rs/zerolog"
)

type Job struct {
	logger zerolog.Logger
}

func NewJob(logger zerolog.Logger) *Job {
	return &Job{logger: logger}
}

func (j *Job) GetOrdersPayment(rpc rpc.RpcInterface, paymentService usecase.Service) background.JobFunc {
	return func() {
		rpc.NewRpcQueue("get_orders_payment", func(data string) string {
			j.logger.Debug().Msg("Get Orders Payment...")
			var request model.GetOrdersPaymentRequest
			bindErr := json.Unmarshal([]byte(data), &request)
			if bindErr != nil {
				j.logger.Error().Err(bindErr).Msg("")
				return ""
			}
			payments := paymentService.GetPaymentsByOrderIds(request.OrderIds)
			ordersPayments := make([]*model.OrderPayment, len(payments))

			j.logger.Debug().Msg("Mapping orders payments result")
			for _, payment := range payments {
				ordersPayments = append(ordersPayments, model.CreateOrderPayment(*payment))
			}
			res := model.OrderPaymentResponse{
				PaymentResult: ordersPayments,
			}
			jsonResponse, bindErr := json.Marshal(res)
			if bindErr != nil {
				j.logger.Error().Err(bindErr).Msg("")
				return ""
			}
			return string(jsonResponse)
		})
	}
}
