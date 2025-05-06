package auth

import (
	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/google/uuid"
)

type LoginResponse struct {
	Id           uuid.UUID `json:"id"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	Username     string    `json:"username"`
	AvatarUrl    string    `json:"avatarUrl"`
}

func From(u entity.User, accessToken, refreshToken string) *LoginResponse {
	return &LoginResponse{
		Id:           u.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     u.Username,
		AvatarUrl:    u.AvatarUrl,
	}
}
