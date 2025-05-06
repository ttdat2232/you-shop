package middleware

import (
	"net/http"

	"github.com/TechwizsonORG/product-service/api/model"
	appErr "github.com/TechwizsonORG/product-service/err"
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
		for _, err := range c.Errors {
			logger.Error().Err(err).Msg("")
		}

		for _, err := range c.Errors {
			logger.Error().Err(err).Msg("Message occurred")
			switch e := err.Err.(type) {
			case appErr.ApplicationError:
				c.JSON(e.Code(), model.NewApiResponse(e.Code(), e.Title(), false, model.ErrorResponse{
					Type:     c.Request.Host + c.Request.URL.Path,
					Title:    e.Title(),
					Detail:   e.Detail(),
					Status:   e.Code(),
					Instance: c.Request.URL.Path,
					Data:     e.Data(),
				}))
				return
			default:
				c.JSON(http.StatusInternalServerError, model.NewApiResponse(http.StatusInternalServerError, "Internal Server Error", false, model.ErrorResponse{
					Type:     c.Request.Host,
					Title:    "Internal Server Error",
					Status:   http.StatusInternalServerError,
					Instance: c.Request.URL.Path,
				}))
				return
			}

		}
	}
}
