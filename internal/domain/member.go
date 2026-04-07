package domain

import (
	"time"
	"github.com/google/uuid"
)

type Member struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name                  string    `gorm:"not null" json:"name" validate:"required"`
	Phone                 string    `gorm:"unique;not null" json:"phone" validate:"required"`
	Status                string    `gorm:"default:inactive" json:"status"`
	FingerprintID 		  *int      `gorm:"uniqueIndex" json:"fingerprint_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}