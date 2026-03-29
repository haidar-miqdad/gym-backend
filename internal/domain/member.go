package domain

import (
	"time"
	"github.com/google/uuid"
)

type Member struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name                  string    `gorm:"not null" json:"name" validate:"required"`
	Phone                 string    `gorm:"unique;not null" json:"phone" validate:"required"`
	FingerprintMappingID  int       `gorm:"unique;index" json:"fingerprint_id"`
	Status                string    `gorm:"default:inactive" json:"status"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}