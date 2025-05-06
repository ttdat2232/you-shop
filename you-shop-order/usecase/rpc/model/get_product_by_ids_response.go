package model

import (
	"github.com/google/uuid"
)

type GetProductByIdsResponse struct {
	Products []Product `json:"products"`
}

type Product struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
