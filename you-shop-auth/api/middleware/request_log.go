package middleware

import (
	"time"

	"github.com/TechwizsonORG/auth-service/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func RequestLog(log zerolog.Logger) gin.HandlerFunc {
	logger := log.
		With().
		Str("Middleware", "Request log").
		Logger()
	return func(c *gin.Context) {
		start := util.GetCurrentUtcTime(7)
		requestIdHeader := "X-Request-Id"
		c.Request.Header.Add(requestIdHeader, uuid.New().String())
		logger.Info().Msgf("Method: %s, Path: %s, RemoteAddr: %s, ContentLength: %d, Host: %s, Referer: %s, User-Agent: %s, RequestId: %s", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, c.Request.ContentLength, c.Request.Host, c.Request.Referer(), c.Request.Header.Get("User-Agent"), c.Request.Header.Get(requestIdHeader))
		c.Next()
		end := util.GetCurrentUtcTime(7)
		logger.Info().Msgf("Finished Request %s from %v to %v - %dms, status %v", c.Request.Header.Get(requestIdHeader), start, end, time.Since(start).Milliseconds(), c.Writer.Status())
	}
}
