package auth

import (
	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/TechwizsonORG/auth-service/err"
)

type AuthInterface interface {
	Login(email string, username string, password string) (u *entity.User, accessToken string, refreshToken string, e *err.AppError)
	Register(email string, username string, password string, phoneNumber string) (*entity.User, *err.AppError)
}
