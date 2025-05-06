package token

import "github.com/TechwizsonORG/auth-service/usecase/token"

type ValidateTokenRequest struct {
	Token string          `json:"token"`
	Type  token.TokenType `json:"type"`
}
