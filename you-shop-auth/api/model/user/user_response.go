package user

import (
	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/google/uuid"
)

type UserResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Scope    string    `json:"scope"`
}

func From(user entity.User, role, scope string) *UserResponse {
	return &UserResponse{
		Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
		Role:     role,
		Scope:    scope,
	}
}
