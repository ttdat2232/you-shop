package handler

import (
	"github.com/TechwizsonORG/price-service/usecase"
	"github.com/gin-gonic/gin"
)

type PriceHandler struct {
	PriceService usecase.Service
}

func NewPriceHandler(priceService usecase.Service) *PriceHandler {
	return &PriceHandler{PriceService: priceService}
}

func (ph *PriceHandler) RegisterRoutes(router *gin.RouterGroup) {}
