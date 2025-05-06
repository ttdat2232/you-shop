package handler

import (
	"fmt"

	"github.com/TechwizsonORG/image-service/api/middleware"
	"github.com/TechwizsonORG/image-service/api/model"
	imageModel "github.com/TechwizsonORG/image-service/api/model/image"
	"github.com/TechwizsonORG/image-service/err"
	"github.com/TechwizsonORG/image-service/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

type ImageHandler struct {
	imageService usecase.Service
	fileService  usecase.FileService
}

func NewImageHandler(imageService usecase.Service, fileService usecase.FileService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
		fileService:  fileService,
	}
}

func (i *ImageHandler) RegisterRoutes(routes *gin.RouterGroup) {
	imageRoutes := routes.Group("/images")
	imageRoutes.GET("/:id", i.serveImageById)
	imageRoutes.GET("/banner", i.getBanners)
	imageRoutes.POST("/upload", i.upload)
	imageRoutes.POST("/upload/banner", middleware.AuthorizationMiddleware([]string{"admin"}, nil), i.uploadBanners)
	imageRoutes.DELETE("/:id", middleware.AuthorizationMiddleware([]string{"admin"}, nil), i.deleteImage)
	imageRoutes.DELETE("/banner", middleware.AuthorizationMiddleware([]string{"admin"}, nil), i.deleteBanners)
}

// Upload godoc
//
//	@Summary	Upload Image
//	@Tags		images
//	@Accept		multipart/form-data
//	@Param		image_file	formData	file	true	"Image File"
//	@Param		alt			formData	string	true	"Image Alt"
//	@Param		owner_id	formData	string	true	"Owner ID"
//	@Failure	400			{object}	model.ApiResponse{data=err.AppError}
//	@Success	201
//	@Response	default
//	@Header		200	{string}	Location	"/api/v1/images/{id}"
//	@Router		/images/upload [post]
func (i *ImageHandler) upload(c *gin.Context) {
	file, header, parseErr := c.Request.FormFile("image_file")
	if parseErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewValidationErr("", "form didn't include image_file field", nil)})
		c.Errors = append(c.Errors, &gin.Error{Err: parseErr})
		return
	}
	defer file.Close()
	ownerId, parseErr := uuid.Parse(c.Request.FormValue("owner_id"))
	if parseErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewValidationErr("", "error when parse owner_id to uuid type", nil)})
		return
	}

	if nImage, appErr := i.imageService.SaveImage(file, header.Filename, c.Request.FormValue("alt"), ownerId); appErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr})
	} else {
		c.Header("Location", fmt.Sprintf("/api/v1/images/%s", nImage.Id.String()))
		c.Status(201)
	}

}

// ServeImageById godoc
//
//	@Summary	Serve Image
//	@Tags		images
//	@Produce	jpeg
//	@Produce	png
//	@Param		id	path	string	true	"Image ID in UUID format (e.g: 123e4567-e89b-12d3-a456-426614174000)"
//	@Router		/images/{id} [get]
func (i *ImageHandler) serveImageById(c *gin.Context) {
	id, parseErr := uuid.Parse(c.Param("id"))

	if parseErr != nil {
		msg := fmt.Sprintf("Cannot parse id %s", c.Param("id"))
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewAppError(400, msg, msg, nil)})
		return
	}

	image := i.imageService.GetImageById(id)
	imageBytes, getFileErr := i.fileService.GetFile(image.ImageUrl)
	if getFileErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: getFileErr})
		c.Writer.Write(make([]byte, 0))
	} else {
		c.Header("Content-Type", image.ContentType)
		c.Writer.Write(imageBytes)
	}
}

// DeleteImage godoc
//
//	@Summary	Delete Image
//	@Tags		images
//	@Param		id	path		string	true	"Image ID in UUID format (e.g: 123e4567-e89b-12d3-a456-426614174000)"
//	@Failure	400	{object}	model.ApiResponse{data=err.AppError}
//	@Success	202
//	@Router		/images/{id} [delete]
func (i *ImageHandler) deleteImage(c *gin.Context) {
	id, parseErr := uuid.Parse(c.Param("id"))

	if parseErr != nil {
		msg := fmt.Sprintf("Cannot parse id %s", c.Param("id"))
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewAppError(400, msg, msg, nil)})
		return
	}
	err := i.imageService.DeleteImage(id)
	if err != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: err})
		return
	}
	c.Status(202)
}

// UploadBanners godoc
//
//	@Summary	Upload Banners
//	@Tags		images
//	@Accept		multipart/form-data
//	@Param		banner_files	formData	file	true	"Banner Files"
//	@Failure	400				{object}	model.ApiResponse{data=err.AppError}
//	@Success	201
//	@Router		/images/upload/banner [post]
func (i *ImageHandler) uploadBanners(c *gin.Context) {

	form, _ := c.MultipartForm()
	files := form.File["banner_files"]
	if ok, saveError := i.imageService.SaveBanners(files); saveError != nil || !ok {
		c.Errors = append(c.Errors, &gin.Error{Err: saveError})
	} else {
		c.Status(201)
	}
}

// DeleteBanners godoc
//
//	@Summary	Delete Banners
//	@Tags		images
//	@Failure	400	{object}	model.ApiResponse{data=err.AppError}
//	@Success	202
//	@Router		/images/banner [delete]
func (i *ImageHandler) deleteBanners(c *gin.Context) {
	ok, deleteErr := i.imageService.DeleteBanners()
	if deleteErr != nil || !ok {
		c.Errors = append(c.Errors, &gin.Error{Err: deleteErr})
	} else {
		c.Status(202)
	}
}

// GetBanner godoc
//
//	@Summary	Get Banners Data
//	@Tags		images
//	@Produce	json
//	@Failure	400	{object}	model.ApiResponse{data=err.AppError}
//	@Success	200	{object}	model.ApiResponse{data=[]imageModel.ImageResponse}
//	@Router		/images/banner [get]
func (i *ImageHandler) getBanners(c *gin.Context) {
	banners := i.imageService.GetBanners()
	imageResponse := []imageModel.ImageResponse{}
	for _, banner := range banners {
		imageResponse = append(imageResponse, imageModel.From(banner, c))
	}
	c.JSON(200, model.SuccessResponse(imageResponse))
}
