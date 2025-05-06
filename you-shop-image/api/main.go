package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TechwizsonORG/image-service/api/docs"
	"github.com/TechwizsonORG/image-service/api/handler"
	"github.com/TechwizsonORG/image-service/api/middleware"
	"github.com/TechwizsonORG/image-service/api/model"
	"github.com/TechwizsonORG/image-service/background"
	"github.com/TechwizsonORG/image-service/config"
	infraFile "github.com/TechwizsonORG/image-service/infrastructure/file"
	"github.com/TechwizsonORG/image-service/infrastructure/repository"
	"github.com/TechwizsonORG/image-service/infrastructure/rpc"
	"github.com/TechwizsonORG/image-service/job"
	"github.com/TechwizsonORG/image-service/usecase"
	"github.com/TechwizsonORG/image-service/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	dbConfig, serverConfig, logConfig, rabbitMqConfig, s3ProxyConfig, httpEndpoint, mode := config.Init()

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

	logger.Info().Msg("Starting server...")
	db, err := sql.Open("postgres", dbConfig.GetPostgresDSN())
	if err != nil {
		logger.Err(err).Msg("Error occurred")
	}
	defer db.Close()

	imageRepo := repository.NewImageRepository(logger, db)
	fileService := infraFile.NewTypeService(*s3ProxyConfig)
	imageService := usecase.NewImageService(imageRepo, fileService)
	imageHandler := handler.NewImageHandler(imageService, fileService)
	rpcService := rpc.NewRpcService(*rabbitMqConfig, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job := job.NewJob(logger)
	background.Go(logger, job.GetOwnersImages(ctx, rpcService, imageService))

	docs.SwaggerInfo.Title = "Image Service API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Information of You Shop Image API Endpoints"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "api.youshop.fun"
	docs.SwaggerInfo.Schemes = []string{"https"}

	gin.SetMode(mode)
	router := gin.New()

	router.Use(middleware.Cors())
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.RequestLog(logger))
	router.Use(middleware.AuthenticateMiddleware(logger, *httpEndpoint))

	v1 := router.Group("/api/v1")
	imageHandler.RegisterRoutes(v1)
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, model.SuccessResponse(nil))
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
}
