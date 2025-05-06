package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditEntity struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
