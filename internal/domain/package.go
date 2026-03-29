package domain

import (
	"time"
	"github.com/google/uuid"
)

type Package struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`         // Contoh: "Daily Pass", "Monthly Pro"
	DurationDays int       `gorm:"not null" json:"duration_days"` // Harian = 1, Bulanan = 30
	Price        float64   `gorm:"not null" json:"price"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}