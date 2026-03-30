package domain

import (
	"time"
	"github.com/google/uuid"
)

type Payment struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SubscriptionID  uuid.UUID `gorm:"type:uuid;index;not null" json:"subscription_id"`
	Amount          float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Method          string    `gorm:"type:varchar(20);not null" json:"method"` // CASH, QRIS, TRANSFER
	ReferenceNumber string    `gorm:"type:varchar(100)" json:"reference_number"`
	Status          string    `gorm:"type:varchar(20);default:'completed'" json:"status"` // COMPLETED, REFUNDED
	ProcessedBy     uuid.UUID `gorm:"type:uuid" json:"processed_by"` // ID Admin/Kasir
	CreatedAt       time.Time `json:"created_at"`
}