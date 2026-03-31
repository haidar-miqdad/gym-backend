package domain

import (
	"time"
	"github.com/google/uuid"
)

type AccessLog struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MemberID       uuid.UUID `gorm:"type:uuid;index;not null" json:"member_id"`
	SubscriptionID uuid.UUID `gorm:"type:uuid;index;not null" json:"subscription_id"`
	CheckInAt      time.Time `json:"check_in_at"`
}