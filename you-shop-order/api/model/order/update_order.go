package order

import (
	"github.com/TechwizsonORG/order-service/entity"
	"github.com/google/uuid"
)

type UpdateOrder struct {
	Description string `json:"description"`
	IsCancel    bool   `json:"isCancel"`
}

func (o UpdateOrder) ToEntity(id uuid.UUID) entity.Order {
	return entity.Order{
		AuditEntity: entity.AuditEntity{
			Id: id,
		},
		Description: o.Description,
	}
}
