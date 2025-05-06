package color

import (
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/google/uuid"
)

type ColorResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func FromColorEntity(color entity.Color) *ColorResponse {
	return &ColorResponse{
		Id:   color.Id,
		Name: color.Name,
	}
}
