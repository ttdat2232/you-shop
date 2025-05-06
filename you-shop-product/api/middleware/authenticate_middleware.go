package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	apiModel "github.com/TechwizsonORG/product-service/api/model"
	"github.com/TechwizsonORG/product-service/api/model/token"
	"github.com/TechwizsonORG/product-service/config/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func AuthenticateMiddleware(logger zerolog.Logger, httpEndpoint model.HttpEndpoint) gin.HandlerFunc {
	logger = logger.
		With().
		Str("Middleware", "Auth").
		Logger()
	return func(c *gin.Context) {
		msg := "Unauthorized request accessed"
		tokenString := c.Request.Header.Get("Authorization")
		tokenReq := token.ValidateTokenRequest{
			Token: tokenString,
			Type:  token.AccessToken,
		}
		jsonBody, _ := json.Marshal(tokenReq)
		bodyReader := bytes.NewReader(jsonBody)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/token/validate", httpEndpoint.AuthServerUrl), bodyReader)
		if err != nil {
			logger.Debug().Err(err).Msg("")
			logger.Warn().Msg(msg)
		} else {
			req.Header.Set("Content-Type", "application/json")
			client := http.Client{
				Timeout: 30 * time.Second,
			}
			res, err := client.Do(req)
			if err != nil {
				logger.Debug().Err(err).Msg("")
				logger.Warn().Msg(msg)
			} else if res != nil && res.StatusCode == 200 {
				defer res.Body.Close()
				setAuthHeader(c, res, &logger)
			} else {
				logger.Warn().Msg(msg)
			}
		}
		c.Next()
	}
}

func setAuthHeader(c *gin.Context, res *http.Response, logger *zerolog.Logger) {
	// This header will be marker to know user have been authenticated
	c.Request.Header.Set("Authenticated", "")

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read response body")
		return
	}

	var apiRes apiModel.ApiResponse
	if err := json.Unmarshal(bodyBytes, &apiRes); err != nil {
		logger.Error().Err(err).Msg("Failed to parse API response")
		return
	}
	dataMap, ok := apiRes.Data.(map[string]interface{})

	if !ok {
		logger.Error().Msg("Failed to assert API response data as map[string]interface{}")
	}

	userId, ok := dataMap["id"].(string)
	if !ok {
		logger.Error().Msg("Failed to extract user ID from API response data")
	} else {
		c.Request.Header.Set("userId", userId)
	}

	role, ok := dataMap["role"].(string)
	if !ok {
		logger.Error().Msg("Failed to extract role from API response data")
	} else {
		c.Request.Header.Set("role", role)
	}

	scope, ok := dataMap["scope"].(string)
	if !ok {
		logger.Error().Msg("Failed to extract scope from API response data")
	} else {
		c.Request.Header.Set("scope", scope)
	}
}
