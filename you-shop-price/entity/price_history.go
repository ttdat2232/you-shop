package entity

import (
	"time"

	"github.com/google/uuid"
)

type PriceHistory struct {
	AuditEntity
	ProductId     uuid.UUID
	OldPriceId    uuid.UUID
	PreviousPrice float64
	NewPrice      float64
	ChangedAt     time.Time
	ChangeBy      time.Time
	ChangeReason  string
}
