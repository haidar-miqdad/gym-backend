package repository

import (
	"context"
	"gym-backend/internal/domain"
	"time"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	FindActiveByMemberID(ctx context.Context, memberID string) (domain.Subscription, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db}
}

func (r *subscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *subscriptionRepository) FindActiveByMemberID(ctx context.Context, memberID string) (domain.Subscription, error) {
	var sub domain.Subscription
	// Logika: Cari sub yang EndDate-nya masih di masa depan
	err := r.db.WithContext(ctx).
		Where("member_id = ? AND end_date >= ?", memberID, time.Now()).
		Order("end_date DESC").
		First(&sub).Error
	return sub, err
}

func (r *subscriptionRepository) GetActiveSubscription(ctx context.Context, memberID string) (domain.Subscription, error) {
	var sub domain.Subscription
	now := time.Now()

	// Logika: Mencari subscription yang rentang waktunya mencakup detik ini
	err := r.db.WithContext(ctx).
		Where("member_id = ? AND start_date <= ? AND end_date >= ? AND status = ?", 
			memberID, now, now, "active").
		First(&sub).Error

	return sub, err
}