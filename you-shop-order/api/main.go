package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TechwizsonORG/order-service/api/docs"
	"github.com/TechwizsonORG/order-service/api/handler"
	"github.com/TechwizsonORG/order-service/api/middleware"
	"github.com/TechwizsonORG/order-service/background"
	"github.com/TechwizsonORG/order-service/config"
	"github.com/TechwizsonORG/order-service/infrastructure/rabbitmq"
	"github.com/TechwizsonORG/order-service/infrastructure/repository"
	"github.com/TechwizsonORG/order-service/infrastructure/rpc"
	"github.com/TechwizsonORG/order-service/job"
	"github.com/TechwizsonORG/order-service/usecase/order"
	"github.com/TechwizsonORG/order-service/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	docs.SwaggerInfo.Title = "Order API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Information of You Shop Order API Endpoints"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "api.youshop.fun"
	docs.SwaggerInfo.Schemes = []string{"https"}

	dbConfig, serverConfig, logConfig, rabbitConfig, rpcEndpoint, httpEndpoint, mode := config.Init()

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
	logger.Info().Msg("Setting up Order Service")

	// infra
	rpcService := rpc.NewRpcService(*rabbitConfig, logger)
	msq := rabbitmq.NewDefaultMessageQueue(*rabbitConfig, logger)
	repo := repository.NewOrderRepository(db, logger)

	// usecase
	orderService := order.NewOrderService(rpcService, *rpcEndpoint, repo, msq, logger)

	// handler
	orderHandler := handler.NewOrderHandler(orderService, logger, rpcService, *rpcEndpoint)
	adminOrderHandler := handler.NewAdminOrderHandler(logger, orderService)

	// background job
	job := job.NewJob(logger)
	background.Go(logger, job.CreateOrder(rpcService, orderService))
	background.Go(logger, job.HandlePaymentStatusChangedEvent(msq, orderService))

	gin.SetMode(mode)
	r := gin.New()

	r.Use(middleware.Cors())
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.AuthenticateMiddleware(logger, *httpEndpoint))
	r.Use(middleware.ErrorHandler(logger))
	r.Use(middleware.RequestLog(logger))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	route := r.Group("/api/v1")
	orderHandler.AddRoute(route)
	adminOrderHandler.AdminOrderRoute(route)

	go logger.Info().Msg("Order Service Started")
	r.Run(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
}
