package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TechwizsonORG/auth-service/api/docs"
	"github.com/TechwizsonORG/auth-service/api/handler"
	"github.com/TechwizsonORG/auth-service/api/middleware"
	"github.com/TechwizsonORG/auth-service/config"
	configModel "github.com/TechwizsonORG/auth-service/config/model"
	"github.com/TechwizsonORG/auth-service/infrastructure/repository"
	"github.com/TechwizsonORG/auth-service/usecase/auth"
	"github.com/TechwizsonORG/auth-service/usecase/token"
	"github.com/TechwizsonORG/auth-service/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	databaseConfig, svrConfig, logConfig, jwtConfig, mode := config.Init()
	logger := createLogger(logConfig)

	db, err := sql.Open("postgres", databaseConfig.GetPostgresDSN())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db, logger)
	roleRepo := repository.NewRoleRepository(db, logger)
	scopeRepo := repository.NewScopeRepository(db, logger)
	tokenService := token.NewTokenService(logger, userRepo, *jwtConfig, roleRepo, scopeRepo)
	authService := auth.NewAuthService(logger, userRepo, tokenService)

	authHandler := handler.NewAuthHandler(authService)
	tokenHandler := handler.NewTokenHandler(tokenService)

	docs.SwaggerInfo.Title = "Auth API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Information of You Shop Auth API Endpoints"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "api.youshop.fun"
	docs.SwaggerInfo.Schemes = []string{"https"}

	gin.SetMode(mode)
	router := gin.New()
	logger.Info().Msg("Starting Auth Service...")

	router.Use(middleware.Cors())
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestLog(logger))
	router.Use(middleware.ErrorHandler(logger))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/api/v1/auth/health", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{
			"message": "Ok",
		})
	})

	v1 := router.Group("/api/v1/auth")

	authHandler.AuthRoutes(v1)
	tokenHandler.TokenRoutes(v1)

	logger.Info().Msgf("Auth Service is running on %s:%d", svrConfig.Host, svrConfig.Port)
	router.Run(fmt.Sprintf("%s:%d", svrConfig.Host, svrConfig.Port))
}

func createLogger(logConfig *configModel.LogConfig) zerolog.Logger {
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
	return logger
}
