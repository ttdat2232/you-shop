package middleware

import (
	"net/http"

	"github.com/TechwizsonORG/auth-service/api/model"
	"github.com/TechwizsonORG/auth-service/err"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func ErrorHandler(log zerolog.Logger) gin.HandlerFunc {
	logger := log.
		With().
		Str("Middleware", "Error Handler").
		Logger()
	return func(c *gin.Context) {
		c.Next()
		for _, e := range c.Errors {
			logger.Error().Err(e.Err).Msg("Error occurred")
		}
		for _, e := range c.Errors {
			logger.Err(e).Msg("")
			switch e := e.Err.(type) {
			case *err.AppError:
				c.JSON(e.Code, model.NewApiResponse(e.Code, e.Title, false, model.ErrorResponse{
					Type:     c.Request.Host + c.Request.URL.Path,
					Title:    e.Title,
					Detail:   e.Message,
					Status:   e.Code,
					Instance: c.Request.URL.Path,
					Data:     e.Data,
				}))
				return
			default:
				c.JSON(http.StatusInternalServerError, model.NewApiResponse(http.StatusInternalServerError, "Internal Server Error", false, model.ErrorResponse{
					Type:     c.Request.Host + c.Request.URL.Path,
					Title:    "Internal Server Error",
					Status:   http.StatusInternalServerError,
					Instance: c.Request.URL.Path,
				}))
				return
			}

		}
	}
}
