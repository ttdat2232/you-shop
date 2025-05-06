package entity

import "github.com/google/uuid"

type ImageType int8

const (
	ProductImage ImageType = iota + 1
	BannerImage
)

type Image struct {
	AuditEntity
	ImageUrl    string
	Filename    string
	ContentType string
	OwnerId     uuid.UUID
	Size        float32
	Width       int
	Height      int
	IsPublic    bool
	Alt         string
	Type        ImageType
}
