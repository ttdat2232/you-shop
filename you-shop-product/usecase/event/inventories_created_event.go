package event

import "github.com/TechwizsonORG/product-service/usecase/inventory/model"

type CreatedInventoriesEvent struct {
	CreatedInventories []model.CreateInventory
}
