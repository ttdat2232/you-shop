package entity

import (
	"strings"

	"github.com/TechwizsonORG/auth-service/err"
	"github.com/TechwizsonORG/auth-service/util"
	"github.com/google/uuid"
)

type User struct {
	AuditEntity
	Username           string
	NormalizedUsername string
	PasswordHash       string
	Email              string
	NormalizedEmail    string
	AvatarUrl          string
	IsActive           bool
	PhoneNumber        string
}

func CreateUser(username string, password string, email string, avatarUrl string, phoneNumber string) (*User, *err.AppError) {
	hashPass, error := util.HashPassword(password)
	if error != nil {
		return nil, err.NewCreateUserError("Hash password failed", nil)
	}
	return &User{
		AuditEntity: AuditEntity{
			Id: uuid.New(),
		},
		Username:           username,
		NormalizedUsername: strings.ToUpper(username),
		PasswordHash:       hashPass,
		Email:              email,
		NormalizedEmail:    strings.ToUpper(email),
		AvatarUrl:          avatarUrl,
		IsActive:           true,
		PhoneNumber:        phoneNumber,
	}, nil
}
