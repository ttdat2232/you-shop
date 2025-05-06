package model

import "github.com/google/uuid"

type GetProductByIdsRequest struct {
	ProductIds []uuid.UUID `json:"productIds"`
}
