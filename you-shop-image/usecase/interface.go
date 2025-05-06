package usecase

import (
	"mime/multipart"

	"github.com/TechwizsonORG/image-service/entity"
	"github.com/TechwizsonORG/image-service/err"
	"github.com/google/uuid"
)

type ImageExtension string

const (
	JPEG ImageExtension = "jpeg"
	JPG  ImageExtension = "jpg"
	PNG  ImageExtension = "png"
	WEBP ImageExtension = "webp"
	GIF  ImageExtension = "gif"
	AVIF ImageExtension = "avif"
)

type FileService interface {
	SaveToRemoteServer(filename string, extension string, file multipart.File) (string, error)
	GetFile(filepath string) ([]byte, *err.AppError)
	DeleteFile(filepath string) *err.AppError
}

type Service interface {
	GetImageById(id uuid.UUID) *entity.Image
	GetImageByOwnerIds(ids []uuid.UUID, scheme, host, path string) map[string][]string
	GetBanners() []entity.Image
	SaveImage(file multipart.File, filename string, alt string, ownerId uuid.UUID) (*entity.Image, *err.AppError)
	SaveBanners([]*multipart.FileHeader) (bool, *err.AppError)
	DeleteImage(id uuid.UUID) *err.AppError
	DeleteBanners() (bool, *err.AppError)
}

type Repository interface {
	GetById(id uuid.UUID) (*entity.Image, error)
	GetByOwnerIds(ownerIds []uuid.UUID) ([]entity.Image, error)
	AddImage(image entity.Image) (int, error)
	UpdateImage(image entity.Image) (int, error)
	DeleteImage(id uuid.UUID) (int, error)
	AddImages([]entity.Image) (int, error)
	DeleteImageByIds(ids []uuid.UUID) (int, error)
	GetByType(imageType entity.ImageType) ([]entity.Image, error)
}
