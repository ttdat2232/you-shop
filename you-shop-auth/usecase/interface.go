package usecase

import (
	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/google/uuid"
)

// User section

type UserReader interface {
	GetUserByID(id uuid.UUID) (*entity.User, error)
	GetUserByEmailOrUsername(email string, username string) (*entity.User, error)
}

type UserWriter interface {
	AddUser(user entity.User) (entity.User, error)
	UpdateUser(user entity.User) (entity.User, error)
}

type UserRepository interface {
	UserReader
	UserWriter
}

// Role section
type RoleReader interface {
	GetRolesByUserId(userId uuid.UUID) []entity.Role
}

type RoleRepository interface {
	RoleReader
}

// Scope section
type ScopeReader interface {
	GetScopesByUserId(userId uuid.UUID) []entity.Scope
}

type ScopeRepository interface {
	ScopeReader
}

// Client section

type ClientReader interface {
	GetClientByID(id uuid.UUID) (entity.Client, error)
}

type ClientWriter interface {
	AddClient(client entity.Client) (entity.Client, error)
	DeleteClient(id uuid.UUID) error
}

type ClientRepository interface {
	ClientReader
	ClientWriter
}

// Authorization Code section

type AuthorizationCodeReader interface {
	GetAuthorizationCode(code string) (entity.AuthorizationCode, error)
}

type AuthorizationCodeWriter interface {
	AddAuthorizationCode(code entity.AuthorizationCode) (entity.AuthorizationCode, error)
}

type AuthorizationCodeRepository interface {
	AuthorizationCodeReader
	AuthorizationCodeWriter
}
