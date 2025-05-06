package handler

import (
	"github.com/TechwizsonORG/product-service/api/model"
	sizeModel "github.com/TechwizsonORG/product-service/api/model/size"
	"github.com/TechwizsonORG/product-service/err"
	"github.com/TechwizsonORG/product-service/usecase/size"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type SizeHandler struct {
	logger      zerolog.Logger
	sizeService size.SizeUseCase
}

func NewSizeHandler(sizeService size.SizeUseCase, logger zerolog.Logger) *SizeHandler {
	logger = logger.With().Str("handler", "size").Logger()
	return &SizeHandler{
		logger:      logger,
		sizeService: sizeService,
	}
}

func (s *SizeHandler) SizeRoute(router *gin.RouterGroup) {
	sizeGroup := router.Group("/sizes")
	sizeGroup.GET("", s.getSizes)
	sizeGroup.POST("", s.createSize)
}

func (s *SizeHandler) getSizes(c *gin.Context) {
	sizes := s.sizeService.GetSizes()
	result := make([]sizeModel.SizeResponse, 0, len(sizes))
	for _, size := range sizes {
		result = append(result, *sizeModel.FromSizeEntity(size))
	}
	c.JSON(200, model.SuccessResponse(result))
}
func (s *SizeHandler) createSize(c *gin.Context) {
	var createSize sizeModel.CreateSizeRequest
	bindErr := c.BindJSON(&createSize)
	if bindErr != nil {
		s.logger.Error().Err(bindErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewProductError(400, "couldn't parse request body", "couldn't parse request body", nil)})
		return
	}
	addedSize, addErr := s.sizeService.AddSize(createSize.Name)
	if addErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: addErr})
		return
	}
	c.JSON(200, model.SuccessResponse(sizeModel.FromSizeEntity(*addedSize)))
}
