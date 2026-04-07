package repository

import (
	"context"
	"gym-backend/internal/domain"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(ctx context.Context, tx *gorm.DB, payment *domain.Payment) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Create(ctx context.Context, tx *gorm.DB, payment *domain.Payment) error {
	// Jika tx disediakan (dari service), gunakan tx tersebut. Jika tidak, gunakan r.db default.
	db := tx
    if db == nil { db = r.db }
    return db.WithContext(ctx).Create(payment).Error
}
