package token

import (
	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/TechwizsonORG/auth-service/err"
	"github.com/google/uuid"
)

type TokenType int8

const (
	AccessToken TokenType = iota
	RefreshToken
)

type TokenInterface interface {
	GenerateToken(userId uuid.UUID, tokenType TokenType) (string, *err.AppError)
	ValidateTokenWithResponse(token string, tokenType TokenType) (user *entity.User, role, scope string, appErr *err.AppError)
}
