package handler

import (
	"github.com/TechwizsonORG/product-service/api/middleware"
	"github.com/TechwizsonORG/product-service/api/model"
	colorModel "github.com/TechwizsonORG/product-service/api/model/color"
	"github.com/TechwizsonORG/product-service/err"
	"github.com/TechwizsonORG/product-service/usecase/color"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ColorHandler struct {
	colorService color.ColorUsecase
	logger       zerolog.Logger
}

func NewColorHandler(logger zerolog.Logger, colorService color.ColorUsecase) *ColorHandler {
	logger = logger.With().Str("handler", "color").Logger()
	return &ColorHandler{
		logger:       logger,
		colorService: colorService,
	}
}

func (c *ColorHandler) ColorRoute(routeGroup *gin.RouterGroup) {
	colorRoute := routeGroup.Group("/colors")
	colorRoute.GET("", middleware.AuthorizationMiddleware([]string{"admin"}, nil), c.getColor)
	colorRoute.GET(":id", middleware.AuthorizationMiddleware([]string{"admin"}, nil), c.getColorById)
	colorRoute.POST("", middleware.AuthorizationMiddleware([]string{"admin"}, nil), c.addColor)
}

func (color *ColorHandler) getColorById(c *gin.Context) {
	colorId, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		color.logger.Error().Err(parseErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewProductError(400, "counldn't parse color id", "couldn't parse color id", nil)})
		return
	}
	colorEntity, getErr := color.colorService.GetById(colorId)
	if getErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: getErr})
		return
	}
	c.JSON(200, model.SuccessResponse(colorModel.FromColorEntity(colorEntity)))

}
func (color *ColorHandler) getColor(c *gin.Context) {
	colors := color.colorService.GetColors()
	result := make([]colorModel.ColorResponse, 0, len(colors))
	for _, colorEntity := range colors {
		result = append(result, *colorModel.FromColorEntity(colorEntity))
	}
	c.JSON(200, model.SuccessResponse(result))
}

func (color *ColorHandler) addColor(c *gin.Context) {
	var createColor colorModel.CreateColorRequest
	bindErr := c.BindJSON(&createColor)
	if bindErr != nil {
		color.logger.Error().Err(bindErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: err.NewProductError(400, "couldn't parse request json", "couldn't parse request json", nil)})
		return
	}
	addedColor, addErr := color.colorService.AddColor(createColor.Name)
	if addErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: addErr})
		return
	}
	c.JSON(200, model.SuccessResponse(colorModel.FromColorEntity(*addedColor)))
}
