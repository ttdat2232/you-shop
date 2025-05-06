package middleware

import (
	"net/http"

	"github.com/TechwizsonORG/order-service/api/model"
	appErr "github.com/TechwizsonORG/order-service/err"
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
				c.JSON(e.ErrCode(), model.NewApiResponse(e.ErrCode(), e.ErrTitle(), false, model.ErrorResponse{
					Type:     c.Request.Host + c.Request.URL.Path,
					Title:    e.ErrTitle(),
					Detail:   e.ErrDetail(),
					Status:   e.ErrCode(),
					Instance: c.Request.URL.Path,
					Data:     e.ErrData(),
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
