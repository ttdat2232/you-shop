package inventory

import (
	appErr "github.com/TechwizsonORG/product-service/err"
	"github.com/TechwizsonORG/product-service/usecase/inventory/model"
	"github.com/google/uuid"
)

type InventoryUseCase interface {
	AddInventories([]model.CreateInventory) appErr.ApplicationError
	UpdateInventory(productId, colorId, sizeId uuid.UUID, quantity int, price float64) appErr.ApplicationError
	ChangeQuantity(productId, colorId, sizeId uuid.UUID, changeAmount int) appErr.ApplicationError
}
