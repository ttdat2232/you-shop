package job

import (
	"context"
	"encoding/json"

	"github.com/TechwizsonORG/price-service/background"
	"github.com/TechwizsonORG/price-service/usecase"
	"github.com/TechwizsonORG/price-service/usecase/event"
	messagequeue "github.com/TechwizsonORG/price-service/usecase/message_queue"
	"github.com/TechwizsonORG/price-service/usecase/rpc"
	rpcModel "github.com/TechwizsonORG/price-service/usecase/rpc/model"
	"github.com/rs/zerolog"
)

type Job struct {
	logger zerolog.Logger
}

func NewJob(logger zerolog.Logger) *Job {
	return &Job{logger: logger}
}

func (j *Job) InventoriesCreatedHandler(msgQueue messagequeue.MessageQueue, priceService usecase.Service) background.JobFunc {
	return func() {
		msgQueue.Consume(
			*messagequeue.NewDefaultExchangeConfig("you_shop", messagequeue.Topic),
			*messagequeue.NewDefaultQueueConfig("inventory_created_price_consumer", "inventory.created"),
			func(data string) error {
				var productCreatedEvent event.CreatedInventoriesEvent
				unmarshalErr := json.Unmarshal([]byte(data), &productCreatedEvent)
				if unmarshalErr != nil {
					j.logger.Error().Err(unmarshalErr).Msg("")
					return unmarshalErr
				}
				_, createPriceErr := priceService.CreateNewPrices(productCreatedEvent)
				if createPriceErr != nil {
					return createPriceErr
				}
				return nil
			})
	}

}

func (j *Job) GetProductsPrice(ctx context.Context, rpcService rpc.RpcInterface, priceService usecase.Service) background.JobFunc {
	return func() {
		rpcService.NewRpcQueue("get_products_price", func(data string) string {

			var ids []string
			err := json.Unmarshal([]byte(data), &ids)
			if err != nil {
				j.logger.Error().Err(err).Msg("Error unmarshal data:")
				return "[]"
			}

			result, appErr := priceService.GetCurrentPrices(ids)
			if appErr != nil {
				return "[]"
			}

			jsonResult, err := json.Marshal(result)
			if err != nil {
				j.logger.Error().Err(err).Msg("Error marshal data:")
				return "[]"
			}
			return string(jsonResult)
		})
	}
}

func (j *Job) UpdatePrice(ctx context.Context, rpcService rpc.RpcInterface, priceService usecase.Service) background.JobFunc {
	return func() {
		rpcService.NewRpcQueue("update_price", func(data string) string {
			var updateReq rpcModel.UpdatePriceRequest
			err := json.Unmarshal([]byte(data), &updateReq)
			if err != nil {
				j.logger.Error().Err(err).Msg("Error unmarshal data:")
				return "{}"
			}
			_, updateError := priceService.UpdatePrice(updateReq.ProductId, updateReq.ColorId, updateReq.SizeId, updateReq.Price)

			updateRes := &rpcModel.UpdatePriceResponse{
				false,
			}
			defaultRes, _ := json.Marshal(updateReq)

			if updateError != nil {
				return string(defaultRes)
			}
			updateRes.IsUpdated = true
			res, _ := json.Marshal(updateRes)
			return string(res)
		})
	}
}

func (j *Job) GetTotalPrice(rpcService rpc.RpcInterface, priceService usecase.Service) background.JobFunc {
	return func() {
		rpcService.NewRpcQueue("get_total_price", func(data string) string {
			defaultResult := "{}"
			var getTotalPriceReq rpcModel.TotalPriceRequest
			parseJsonErr := json.Unmarshal([]byte(data), &getTotalPriceReq)
			if parseJsonErr != nil {
				j.logger.Error().Err(parseJsonErr).Msg("")
				return defaultResult
			}
			totalPrice, prices, getTotalPriceError := priceService.GetTotalPrice(getTotalPriceReq)
			
			if getTotalPriceError != nil {
				return defaultResult
			}
			res := rpcModel.From(totalPrice, prices, getTotalPriceReq.Items)
			parseResult, parseJsonErr := json.Marshal(res)
			if parseJsonErr != nil {
				return defaultResult
			}
			return string(parseResult)
		})
	}
}
