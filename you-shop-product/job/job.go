package job

import (
	"encoding/json"

	"github.com/TechwizsonORG/product-service/background"
	"github.com/TechwizsonORG/product-service/infrastructure/rpc"
	"github.com/TechwizsonORG/product-service/usecase/inventory"
	messagequeue "github.com/TechwizsonORG/product-service/usecase/message_queue"
	"github.com/TechwizsonORG/product-service/usecase/message_queue/event"
	"github.com/TechwizsonORG/product-service/usecase/product"
	"github.com/TechwizsonORG/product-service/usecase/rpc/model"
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

func (j *Job) CheckProductQuantity(rpc rpc.Service, productService product.UseCase) background.JobFunc {
	return func() {
		rpc.NewRpcQueue("check_product_quantity", func(data string) string {
			var checkProductQuantity model.CheckProductQuantity
			response := model.CheckProductQuantityResponse{
				IsEnough: false,
			}
			defaultRes, _ := json.Marshal(response)
			decodeJsonErr := json.Unmarshal([]byte(data), &checkProductQuantity)
			if decodeJsonErr != nil {
				j.logger.Error().Err(decodeJsonErr).Msg("")
				return string(defaultRes)
			}

			ok, checkErr := productService.CheckProductQuantity(checkProductQuantity.ProductId, checkProductQuantity.ColorId, checkProductQuantity.SizeId, checkProductQuantity.RequireQuantity)
			if checkErr != nil {
				j.logger.Error().Err(checkErr).Msg("")
				return string(defaultRes)
			}
			response.IsEnough = ok
			jsonRes, parseJsonErr := json.Marshal(response)
			if parseJsonErr != nil {
				j.logger.Error().Err(parseJsonErr).Msg("")
				return string(defaultRes)
			}
			return string(jsonRes)
		})
	}
}

func (j *Job) GetProductByIds(rpcService rpc.Service, productService product.UseCase) background.JobFunc {
	return func() {
		rpcService.NewRpcQueue("get_product_by_ids", func(data string) string {
			result := &model.GetProductByIdsResponse{}
			defaultRes, _ := json.Marshal(result)
			var request model.GetProductByIdsRequest
			parseErr := json.Unmarshal([]byte(data), &request)
			if parseErr != nil {
				j.logger.Error().Err(parseErr).Msg("")
				return string(defaultRes)
			}
			products := productService.GetProductByIds(request.ProductIds)
			result = model.From(products)

			res, parseErr := json.Marshal(result)
			if parseErr != nil {
				j.logger.Error().Err(parseErr).Msg("")
				return string(defaultRes)
			}
			return string(res)
		})
	}
}

func (j *Job) OrderUpdatedHandler(msq messagequeue.MessageQueue, inventory inventory.InventoryUseCase) background.JobFunc {
	return func() {
		msq.Consume(
			*messagequeue.NewDefaultExchangeConfig("you_shop", messagequeue.Topic),
			*messagequeue.NewDefaultQueueConfig("product_consumer_updated_order", "order.updated"),
			func(data string) error {
				var updatedOrderEvent event.UpdatedOrderEvent
				json.Unmarshal([]byte(data), &updatedOrderEvent)
				if updatedOrderEvent.Status == event.Confirmed {
					for _ , item := range updatedOrderEvent.Items {
						inventory.ChangeQuantity(item.ProductId, item.ColorId, item.SizeId, -item.Quantity)
					}
				}
				if updatedOrderEvent.Status == event.Returned {
					for _ , item := range updatedOrderEvent.Items {
						inventory.ChangeQuantity(item.ProductId, item.ColorId, item.SizeId, item.Quantity)
					}
				}
				return nil
			},
		)
	}
}
