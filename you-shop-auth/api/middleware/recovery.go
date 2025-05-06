package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Recovery(logger zerolog.Logger) gin.HandlerFunc {
	logger = logger.With().
		Str("Middleware", "Recovery").
		Logger()
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error().Msgf("Panic: %v", err)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
