package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TechwizsonORG/payment-service/api/docs"
	"github.com/TechwizsonORG/payment-service/api/handler"
	"github.com/TechwizsonORG/payment-service/api/middleware"
	"github.com/TechwizsonORG/payment-service/background"
	"github.com/TechwizsonORG/payment-service/config"
	"github.com/TechwizsonORG/payment-service/infrastructure"
	"github.com/TechwizsonORG/payment-service/infrastructure/rabbitmq"
	"github.com/TechwizsonORG/payment-service/infrastructure/repository"
	"github.com/TechwizsonORG/payment-service/infrastructure/rpc"
	"github.com/TechwizsonORG/payment-service/job"
	"github.com/TechwizsonORG/payment-service/usecase"
	"github.com/TechwizsonORG/payment-service/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	docs.SwaggerInfo.Title = "Payment API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Information of You Shop Payment API Endpoints"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "api.youshop.fun"
	docs.SwaggerInfo.Schemes = []string{"https"}

	dbConfig, serverConfig, logConfig, httpEndpoint, rabbitMqConfig, vnpayConfig, rpcEnpoint, mode := config.Init()

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
		panic(err)
	}
	defer db.Close()

	// infra
	rpcService := rpc.NewRpcService(*rabbitMqConfig, logger)
	msq := rabbitmq.NewDefaultMessageQueue(*rabbitMqConfig, logger)
	paymentRepo := repository.NewPaymentRepository(db)

	// service
	vnpayService := infrastructure.NewVnpayService(*vnpayConfig, rpcService, *rpcEnpoint, logger)
	paymentService := usecase.NewPaymentService(paymentRepo, logger, *rpcEnpoint, rpcService, msq)

	// handler
	vnpayHandler := handler.NewVnpayHandler(vnpayService, paymentService, logger)
	paymentHandler := handler.NewPaymentHandler(paymentService, logger)

	//Register job
	jobs := job.NewJob(logger)
	background.Go(logger, jobs.GetOrdersPayment(rpcService, paymentService))

	gin.SetMode(mode)
	route := gin.New()
	route.Use(gin.LoggerWithWriter(logger))
	route.Use(middleware.Cors())
	route.Use(middleware.RequestLog(logger))
	route.Use(middleware.ErrorHandler(logger))
	route.Use(middleware.AuthenticateMiddleware(logger, *httpEndpoint))
	route.GET("/", func(c *gin.Context) {
		c.JSON(200, "hello from payment")
	})

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiGroup := route.Group("/api/v1")
	vnpayHandler.AddVnpayRoute(apiGroup)
	paymentHandler.AddPaymentRoute(apiGroup)

	logger.Info().Msg("Application starting")
	route.Run(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
}
