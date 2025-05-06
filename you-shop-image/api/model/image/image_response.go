package image

import (
	"fmt"

	"github.com/TechwizsonORG/image-service/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageResponse struct {
	Id          uuid.UUID `json:"id"`
	ImageUrl    string    `json:"image_url"`
	ContentType string    `json:"content_type"`
	Alt         string    `json:"alt"`
}

func From(image entity.Image, c *gin.Context) ImageResponse {
	imageUrl := fmt.Sprintf("https://%s/api/v1/images/%s", c.Request.Host, image.Id)
	return ImageResponse{
		Id:          image.Id,
		ImageUrl:    imageUrl,
		ContentType: image.ContentType,
		Alt:         image.Alt,
	}
}
