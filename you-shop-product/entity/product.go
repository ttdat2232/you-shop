package entity

import (
	"time"

	"github.com/TechwizsonORG/product-service/util"
	"github.com/google/uuid"
)

type ProductStatus int64

const (
	Active   ProductStatus = 1
	Inactive ProductStatus = 2
)

type Product struct {
	Id          uuid.UUID
	Name        string
	Description string
	Sku         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Thumbnail   string
	UserManual  string
	Status      ProductStatus
}

func From(id uuid.UUID, name string, description string, sku string, createdAt time.Time, updatedAt time.Time, status ProductStatus, thumbnail string, userManual string) *Product {
	return &Product{
		Id:          id,
		Name:        name,
		Description: description,
		Sku:         sku,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Status:      status,
		Thumbnail:   thumbnail,
		UserManual:  userManual,
	}
}

func NewProduct(name string, description string, sku string, userManual string) *Product {
	current := util.GetCurrentUtcTime(7)
	return &Product{
		Id:          uuid.New(),
		Name:        name,
		Description: description,
		Sku:         sku,
		CreatedAt:   current,
		UpdatedAt:   current,
		Status:      Active,
		UserManual:  userManual,
	}
}

func (p *Product) Update(name string, description string, sku string, status ProductStatus, userManual string) {
	p.Name = name
	p.Description = description
	p.Sku = sku
	p.UpdatedAt = util.GetCurrentUtcTime(7)
	p.Status = status
	p.UpdatedAt = util.GetCurrentUtcTime(7)
	p.UserManual = userManual
}
