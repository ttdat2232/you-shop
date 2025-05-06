package entity

import "github.com/google/uuid"

type ProductImage struct {
	Id        uuid.UUID
	ProductId uuid.UUID
	ColorId   uuid.UUID
	ImageUrl  string
	IsPrimary bool
	IsPublic  bool
}

func NewPrimaryProductImage(id uuid.UUID, imageUrl string) *ProductImage {
	return &ProductImage{
		Id:        id,
		ImageUrl:  imageUrl,
		IsPrimary: true,
		IsPublic:  true,
	}
}
