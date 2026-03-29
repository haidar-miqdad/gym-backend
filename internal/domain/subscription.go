package domain

import (
	"time"
	"github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MemberID  uuid.UUID `gorm:"type:uuid;index" json:"member_id"`
	PackageID uuid.UUID `gorm:"type:uuid" json:"package_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `gorm:"default:active" json:"status"` // active, expired, cancelled
	CreatedAt time.Time `json:"created_at"`
}