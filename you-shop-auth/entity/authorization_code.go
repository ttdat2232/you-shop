package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuthorizationCode struct {
	AuditEntity
	Code        string
	ClientId    uuid.UUID
	UserId      uuid.UUID
	RedirectUri string
	ExpiresAt   time.Time
	Scopes      string
}
