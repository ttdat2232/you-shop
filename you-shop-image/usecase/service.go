package usecase

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/TechwizsonORG/image-service/entity"
	"github.com/TechwizsonORG/image-service/err"
	"github.com/TechwizsonORG/image-service/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ImageService struct {
	imageRepo   Repository
	fileService FileService
	logger      zerolog.Logger
}

func NewImageService(imageRepo Repository, fileService FileService) *ImageService {
	return &ImageService{
		imageRepo:   imageRepo,
		fileService: fileService,
	}
}

func isValidExtension(ext string) bool {
	allowedExtensions := []ImageExtension{
		AVIF,
		GIF,
		JPEG,
		JPG,
		PNG,
		WEBP,
	}
	ext = strings.ToLower(ext)
	for _, allowed := range allowedExtensions {
		if ext == string(allowed) {
			return true
		}
	}
	return false
}

func (i *ImageService) SaveImage(file multipart.File, filename string, alt string, ownerId uuid.UUID) (*entity.Image, *err.AppError) {
	splitFile := strings.Split(filename, ".")
	extension := splitFile[len(splitFile)-1]
	if !isValidExtension(extension) {
		return nil, err.NewAppError(400, "File's extension is not allowed", "File's extension is not allowed", nil)
	}

	imageDst, saveErr := i.fileService.SaveToRemoteServer(fmt.Sprintf("%d%s", util.GetCurrentUtcTime(7).UnixNano(), util.RandString(10)), extension, file)
	if saveErr != nil {
		i.logger.Error().Err(saveErr).Msg("saved to remote server failed")
		return nil, err.NewAppError(500, "Failed when saving image", "Failed when saving image", nil)
	}

	nImage := &entity.Image{
		AuditEntity: entity.AuditEntity{
			Id: uuid.New(),
		},
		ImageUrl:    imageDst,
		Filename:    filename,
		ContentType: getContentType(extension),
		OwnerId:     ownerId,
		IsPublic:    true,
		Alt:         alt,
		Type:        entity.ProductImage,
	}
	_, addErr := i.imageRepo.AddImage(*nImage)
	if addErr != nil {
		i.logger.Error().Err(addErr).Msg("Error Occurred")
		return nil, err.NewAppError(500, "Failed when saving image", "Failed when saving image", nil)
	}
	return nImage, nil
}

func getContentType(ext string) string {
	ext = strings.ToLower(ext)
	var result string
	switch ext {
	case string(PNG):
		result = "image/png"
	case string(JPG):
	case string(JPEG):
		result = "image/jpeg"
	case string(GIF):
		result = "image/gif"
	case string(AVIF):
		result = "image/avif"
	default:
		result = "image/jpeg"
	}
	return result
}

func (i *ImageService) GetImageById(id uuid.UUID) *entity.Image {
	if image, getErr := i.imageRepo.GetById(id); getErr != nil {
		i.logger.Error().Err(getErr).Msg("Fetched image error")
		return &entity.Image{
			Alt:         "Not found",
			ImageUrl:    "https://picsum.photos/500.jpg",
			ContentType: "image/jpg",
		}
	} else {
		return image
	}

}

func (i *ImageService) DeleteImage(id uuid.UUID) *err.AppError {
	image, getErr := i.imageRepo.GetById(id)
	if getErr != nil {
		i.logger.Error().Err(getErr).Msg("Error occurred")
		return err.NewAppError(500, "Failed when deleting image", "Failed when deleting image", nil)
	}
	if fileDeleteErr := i.fileService.DeleteFile(image.ImageUrl); fileDeleteErr != nil {
		i.logger.Error().Err(fileDeleteErr).Msg("Error occurred")
		return err.NewAppError(500, "Failed when deleting image", "Failed when deleting image", nil)
	}
	if _, deleteErr := i.imageRepo.DeleteImage(id); deleteErr != nil {
		i.logger.Error().Err(deleteErr).Msg("Error occurred")
		return err.NewAppError(500, "Failed when deleting image", "Failed when deleting image", nil)
	}
	return nil
}

func (i *ImageService) GetImageByOwnerIds(ids []uuid.UUID, scheme, host, path string) map[string][]string {
	result := make(map[string][]string, len(ids))
	images, getImgErr := i.imageRepo.GetByOwnerIds(ids)
	if getImgErr != nil {
		i.logger.Error().Err(getImgErr).Msg("get image by owner ids failed")
		return result
	}

	for _, image := range images {
		url := fmt.Sprintf("%s://%s/%s/%s", scheme, host, path, image.Id.String())
		result[image.OwnerId.String()] = append(result[image.OwnerId.String()], url)
	}
	return result
}

func (i *ImageService) SaveBanners(fileHeaders []*multipart.FileHeader) (bool, *err.AppError) {

	for _, fileHeader := range fileHeaders {
		splitFile := strings.Split(fileHeader.Filename, ".")
		extension := splitFile[len(splitFile)-1]
		if !isValidExtension(extension) {
			return false, err.NewAppError(400, "File's extension is not allowed", "File's extension is not allowed", nil)
		}
	}

	images := []entity.Image{}
	for _, fileHeader := range fileHeaders {
		file, openErr := fileHeader.Open()
		if openErr != nil {
			i.logger.Error().Err(openErr).Msg("Error occurred")
			return false, err.NewAppError(500, "Failed when saving image", "Failed when saving image", nil)
		}
		defer file.Close()
		imageDst, saveErr := i.fileService.SaveToRemoteServer(fmt.Sprintf("%d%s", util.GetCurrentUtcTime(7).UnixNano(), util.RandString(10)), strings.Split(fileHeader.Filename, ".")[1], file)
		if saveErr != nil {
			i.logger.Error().Err(saveErr).Msg("Error occurred")
			return false, err.NewAppError(500, "Failed when saving image", "Failed when saving image", nil)
		}
		images = append(images, entity.Image{
			AuditEntity: entity.AuditEntity{
				Id: uuid.New(),
			},
			ImageUrl:    imageDst,
			Filename:    fileHeader.Filename,
			ContentType: getContentType(strings.Split(fileHeader.Filename, ".")[1]),
			OwnerId:     uuid.Nil,
			IsPublic:    true,
			Alt:         "Banner",
			Type:        entity.BannerImage,
		})
	}
	_, addImageErr := i.imageRepo.AddImages(images)
	if addImageErr != nil {
		i.logger.Error().Err(addImageErr).Msg("")
		return false, err.NewAppError(500, "Failed when saving image", "Failed when saving image", nil)
	}
	return true, nil
}

func (i *ImageService) DeleteBanners() (bool, *err.AppError) {
	images, getErr := i.imageRepo.GetByType(entity.BannerImage)
	if getErr != nil {
		i.logger.Error().Err(getErr).Msg("Error occurred")
		return false, err.NewAppError(500, "Failed when deleting image", "Failed when deleting image", nil)
	}
	deletedIds := []uuid.UUID{}
	for _, image := range images {
		if fileDeleteErr := i.fileService.DeleteFile(image.ImageUrl); fileDeleteErr != nil {
			i.logger.Error().Err(fileDeleteErr).Msg("Error occurred")
			return false, err.NewAppError(500, "Failed when deleting image", "Failed when deleting image", nil)
		}
		deletedIds = append(deletedIds, image.Id)
	}
	if _, deleteErr := i.imageRepo.DeleteImageByIds(deletedIds); deleteErr != nil {
		i.logger.Error().Err(deleteErr).Msg("Error occurred")
		return false, err.NewAppError(500, "Failed when deleting image", "Failed when deleting image", nil)
	}
	return true, nil
}

func (i *ImageService) GetBanners() []entity.Image {
	images, getErr := i.imageRepo.GetByType(entity.BannerImage)
	if getErr != nil {
		i.logger.Error().Err(getErr).Msg("Error occurred")
		return []entity.Image{}
	}
	return images
}
