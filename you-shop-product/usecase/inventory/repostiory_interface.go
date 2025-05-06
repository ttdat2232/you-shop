package inventory

import (
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/TechwizsonORG/product-service/usecase/inventory/model"
	"github.com/google/uuid"
)

type InventoryRepository interface {
	GetQuantity(productId, sizeId, colorId uuid.UUID) (int, error)
	AddInventories(createInventories []model.CreateInventory) error
	GetInventory(productId uuid.UUID, colorId uuid.UUID, sizeId uuid.UUID) (*entity.Inventory, error)
	UpdateInventory(updateInventory *entity.Inventory) error
}
