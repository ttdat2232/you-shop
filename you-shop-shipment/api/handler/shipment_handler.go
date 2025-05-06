package handler

import (
	"strings"

	apiModel "github.com/TechwizsonORG/shipment-service/api/model"
	ghnService "github.com/TechwizsonORG/shipment-service/infrastructure/ghn/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ShipmentHandler struct {
	logger     zerolog.Logger
	ghnService *ghnService.GhnService
}

func NewShimpentHandler(logger zerolog.Logger, ghnService *ghnService.GhnService) *ShipmentHandler {
	logger = logger.With().Str("handler", "shipment").Logger()
	return &ShipmentHandler{
		logger:     logger,
		ghnService: ghnService,
	}
}

func (s *ShipmentHandler) AddShipmentRoute(group *gin.RouterGroup) {
	shipmentGroup := group.Group("/shipment")
	shipmentGroup.GET("/:provider/provinces", s.getProvinces)
	shipmentGroup.GET("/:provider/districts/", s.getDistricts)
	shipmentGroup.GET("/:provider/wards", s.getWards)
}

func (s *ShipmentHandler) getProvinces(c *gin.Context) {
	provider := c.Param("provider")
	switch {
	case strings.EqualFold(provider, "GHN"):
		provinces := s.ghnService.GetProvinces()
		c.JSON(200, apiModel.SuccessResponse(provinces))
	}
}

func (s *ShipmentHandler) getDistricts(c *gin.Context) {
	provider := c.Param("provider")
	provinceId := c.Query("provinceId")
	switch {
	case strings.EqualFold(provider, "GHN"):

		if provinces, getErr := s.ghnService.GetDistrictsByProvinceId(provinceId); getErr != nil {
			s.logger.Error().Err(getErr).Msg("")
			c.JSON(200, apiModel.SuccessResponse(nil))
		} else {
			c.JSON(200, apiModel.SuccessResponse(provinces))
		}
	}
}
func (s *ShipmentHandler) getWards(c *gin.Context) {
	provider := c.Param("provider")
	districtId := c.Query("districtId")
	switch {
	case strings.EqualFold(provider, "GHN"):
		if wards, getErr := s.ghnService.GetWardsByDisctrictId(districtId); getErr != nil {
			s.logger.Error().Err(getErr).Msg("")
			c.JSON(200, apiModel.SuccessResponse(nil))
		} else {
			c.JSON(200, apiModel.SuccessResponse(wards))
		}
	}
}
