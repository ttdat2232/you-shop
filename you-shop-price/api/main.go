package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TechwizsonORG/price-service/api/handler"
	"github.com/TechwizsonORG/price-service/api/middleware"
	"github.com/TechwizsonORG/price-service/api/model"
	"github.com/TechwizsonORG/price-service/background"
	"github.com/TechwizsonORG/price-service/config"
	"github.com/TechwizsonORG/price-service/infrastructure/rabbitmq"
	"github.com/TechwizsonORG/price-service/infrastructure/repository"
	"github.com/TechwizsonORG/price-service/infrastructure/rpc"
	"github.com/TechwizsonORG/price-service/job"
	"github.com/TechwizsonORG/price-service/usecase"
	"github.com/TechwizsonORG/price-service/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func main() {
	dbConfig, srvConfig, logConfig, _, rabbitMqConfig, mode := config.Init()

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

	db, err := sql.Open("postgres", dbConfig.GetPostgresDSN())
	if err != nil {
		logger.Error().Err(err).Msg("Error occurred")
	}
	defer db.Close()

	rpcService := rpc.NewRpcService(*rabbitMqConfig, logger)
	priceRepo := repository.NewPriceRepository(logger, db)
	priceService := usecase.NewPriceService(logger, priceRepo)
	priceHandler := handler.NewPriceHandler(priceService)
	msgQueue := rabbitmq.NewDefaultMessageQueue(*rabbitMqConfig, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job := job.NewJob(logger)
	background.Go(logger, job.GetProductsPrice(ctx, rpcService, priceService))
	background.Go(logger, job.InventoriesCreatedHandler(msgQueue, priceService))
	background.Go(logger, job.UpdatePrice(ctx, rpcService, priceService))
	background.Go(logger, job.GetTotalPrice(rpcService, priceService))

	logger.Info().Msg("Starting server...")

	gin.SetMode(mode)
	router := gin.New()

	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestLog(logger))
	router.Use(middleware.ErrorHandler(logger))

	v1 := router.Group("/api/v1")

	priceHandler.RegisterRoutes(v1)

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, model.SuccessResponse(nil))
	})

	logger.Info().Msg("Application is running")
	router.Run(fmt.Sprintf("%s:%d", srvConfig.Host, srvConfig.Port))
}
