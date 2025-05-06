package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TechwizsonORG/product-service/api/docs"
	"github.com/TechwizsonORG/product-service/api/handler"
	"github.com/TechwizsonORG/product-service/api/middleware"
	"github.com/TechwizsonORG/product-service/api/model"
	"github.com/TechwizsonORG/product-service/background"
	"github.com/TechwizsonORG/product-service/config"
	"github.com/TechwizsonORG/product-service/infrastructure/rabbitmq"
	"github.com/TechwizsonORG/product-service/infrastructure/repository"
	rpcImpl "github.com/TechwizsonORG/product-service/infrastructure/rpc"
	"github.com/TechwizsonORG/product-service/job"
	"github.com/TechwizsonORG/product-service/usecase/color"
	"github.com/TechwizsonORG/product-service/usecase/inventory"
	"github.com/TechwizsonORG/product-service/usecase/product"
	"github.com/TechwizsonORG/product-service/usecase/size"
	"github.com/TechwizsonORG/product-service/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	docs.SwaggerInfo.Title = "Product API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Information of You Shop Product API Endpoints"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "api.youshop.fun"
	docs.SwaggerInfo.Schemes = []string{"https"}

	dbConfig, srvConfig, logConfig, rabbitMqConfig, rpcServerEndpoint, httpEndpoint, mode := config.Init()

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

	defer func() {
		logger.Info().Msg("Application is shutting down")
	}()
	logger.Info().Msg("Starting server...")
	db, err := sql.Open("postgres", dbConfig.GetPostgresDSN())
	if err != nil {
		logger.Error().Err(err).Msg("Error occurred")
	}
	defer db.Close()

	// infra
	productRepo := repository.NewProductRepository(db, logger)
	inventoryRepo := repository.NewInventoryRepository(db, logger)
	colorRepo := repository.NewColorRepository(db)
	sizeRepo := repository.NewSizeRepository(db)
	msgQueue := rabbitmq.NewDefaultMessageQueue(*rabbitMqConfig, logger)
	rpcService := rpcImpl.NewRpcService(*rabbitMqConfig, logger)

	// service
	productService := product.NewService(*httpEndpoint, productRepo, logger, msgQueue, rpcService, *rpcServerEndpoint, inventoryRepo)
	inventoryService := inventory.NewInventoryService(logger, inventoryRepo, msgQueue, rpcService, *rpcServerEndpoint)
	colorService := color.NewColorService(logger, colorRepo)
	sizeSerivce := size.NewSizeService(sizeRepo, logger)

	// handler
	productHandler := handler.NewProductHandler(productService, rpcService, *rpcServerEndpoint, logger, inventoryService)
	colorHandler := handler.NewColorHandler(logger, colorService)
	sizeHandler := handler.NewSizeHandler(sizeSerivce, logger)

	// job
	job := job.NewJob(logger)
	background.Go(logger, job.CheckProductQuantity(*rpcService, productService))
	background.Go(logger, job.GetProductByIds(*rpcService, productService))
	background.Go(logger, job.OrderUpdatedHandler(msgQueue, inventoryService))

	// gin
	gin.SetMode(mode)
	router := gin.New()

	router.Use(middleware.Cors())
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestLog(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.AuthenticateMiddleware(logger, *httpEndpoint))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")

	v1.GET("/health", func(c *gin.Context) {
		result := make(map[string]string)
		err := db.Ping()
		if err == nil {
			result["database"] = "Healthy"
		} else {
			logger.Error().Err(err).Msg("Error when ping database")
			result["database"] = "Unhealthy"
		}
		c.JSON(200, model.SuccessResponse(result))
	})

	productHandler.ProductRoutes(v1)
	colorHandler.ColorRoute(v1)
	sizeHandler.SizeRoute(v1)

	logger.Info().Msg("Application is running")
	router.Run(fmt.Sprintf("%s:%d", srvConfig.Host, srvConfig.Port))
}
