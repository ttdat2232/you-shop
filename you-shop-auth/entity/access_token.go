package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	AuditEntity
	Token     string
	ClientId  uuid.UUID
	RevokedAt time.Time
}
