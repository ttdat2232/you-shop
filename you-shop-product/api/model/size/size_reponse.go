package size

import (
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/google/uuid"
)

type SizeResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func FromSizeEntity(color entity.Size) *SizeResponse {
	return &SizeResponse{
		Id:   color.Id,
		Name: color.Name,
	}
}
