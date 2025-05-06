package inventory

import (
	"encoding/json"

	configModel "github.com/TechwizsonORG/product-service/config/model"
	"github.com/TechwizsonORG/product-service/err"
	"github.com/TechwizsonORG/product-service/usecase/event"
	"github.com/TechwizsonORG/product-service/usecase/inventory/model"
	messagequeue "github.com/TechwizsonORG/product-service/usecase/message_queue"
	"github.com/TechwizsonORG/product-service/usecase/rpc"
	rpcModel "github.com/TechwizsonORG/product-service/usecase/rpc/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type InventoryService struct {
	inventoryRepo InventoryRepository
	logger        zerolog.Logger
	msq           messagequeue.MessageQueue
	rpcService    rpc.RpcInterface
	rpcEndpoint   configModel.RpcServerEndpoint
}

func NewInventoryService(logger zerolog.Logger, inventoryRepo InventoryRepository, msq messagequeue.MessageQueue, rpcService rpc.RpcInterface, rpcEndpoint configModel.RpcServerEndpoint) *InventoryService {
	logger = logger.With().Str("Inventory", "Service").Logger()
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		logger:        logger,
		msq:           msq,
		rpcService:    rpcService,
		rpcEndpoint:   rpcEndpoint,
	}
}

func (i *InventoryService) AddInventories(createInventories []model.CreateInventory) err.ApplicationError {
	createErr := i.inventoryRepo.AddInventories(createInventories)
	if createErr != nil {
		i.logger.Error().Err(createErr).Msg("")
		return err.NewProductError(500, "creating inventory failed", "creating inventory failed", nil)
	}
	i.msq.Publish(
		*messagequeue.NewDefaultExchangeConfig("you_shop", messagequeue.Topic),
		*messagequeue.NewDefaultQueueConfig("created", "inventory.created"),
		event.CreatedInventoriesEvent{
			CreatedInventories: createInventories,
		},
	)
	return nil
}
func (i *InventoryService) UpdateInventory(productId, colorId, sizeId uuid.UUID, quantity int, price float64) err.ApplicationError {

	updateInventory, getErr := i.inventoryRepo.GetInventory(productId, colorId, sizeId)
	if getErr != nil {
		i.logger.Error().Err(getErr).Msg("")
		return err.CommonError()
	}
	if updateInventory == nil {
		return err.NotFoundProductError("not found inventory")
	}
	jsonReq, _ := json.Marshal(rpcModel.UpdatePriceRequest{
		ProductId: productId,
		ColorId:   colorId,
		SizeId:    sizeId,
		Price:     price,
	})
	result := i.rpcService.Req(i.rpcEndpoint.UpdatePrice, string(jsonReq))
	var res rpcModel.UpdatePriceResponse
	json.Unmarshal([]byte(result), &res)
	if !res.IsUpdated {
		return err.NewProductError(500, "updateing price failed", "updateing price failed", nil)
	}

	updateInventory.Quantity = quantity
	updateErr := i.inventoryRepo.UpdateInventory(updateInventory)
	if updateErr != nil {
		i.logger.Error().Err(updateErr).Msg("")
		return err.CommonError()
	}
	return nil
}

func (i *InventoryService) ChangeQuantity(productId, colorId, sizeId uuid.UUID, changeAmount int) err.ApplicationError {
	inventory, getErr := i.inventoryRepo.GetInventory(productId, colorId, sizeId)
	if getErr != nil {
		i.logger.Error().Err(getErr).Msg("")
		return err.CommonError()
	}
	if inventory == nil {
		return err.NewProductError(404, "counldn't found inventory", "counldn't found inventory", nil)
	}

	inventory.Quantity = inventory.Quantity + changeAmount
	updateErr := i.inventoryRepo.UpdateInventory(inventory)
	if updateErr != nil {
		i.logger.Error().Err(getErr).Msg("")
		return err.CommonError()
	}
	return nil
}
