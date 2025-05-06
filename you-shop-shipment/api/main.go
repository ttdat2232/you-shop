package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TechwizsonORG/shipment-service/api/handler"
	"github.com/TechwizsonORG/shipment-service/config"
	"github.com/TechwizsonORG/shipment-service/infrastructure/ghn/service"
	"github.com/TechwizsonORG/shipment-service/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	_, serverConfig, logConfig, _, _, ghnConfig, mode := config.Init()

	zerolog.TimestampFunc = func() time.Time {
		return util.GetCurrentUtcTime(7)
	}
	multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
	})
	logger := zerolog.New(multi).
		Level(zerolog.Level(logConfig.Level)).
		With().
		Timestamp().
		Caller().
		Logger()

	//infrastructure
	ghnService := service.NewGhnService(*ghnConfig, logger)

	shipmentHandler := handler.NewShimpentHandler(logger, ghnService)
	gin.SetMode(mode)
	route := gin.New()
	v1 := route.Group("/api/v1")
	shipmentHandler.AddShipmentRoute(v1)
	route.Run(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
}
