package repository

import (
	"context"
	"gym-backend/internal/domain"
	"gorm.io/gorm"
)

type AccessLogRepository interface {
	Create(ctx context.Context, log *domain.AccessLog) error
}

type accessLogRepository struct {
	db *gorm.DB
}

func NewAccessLogRepository(db *gorm.DB) AccessLogRepository {
	return &accessLogRepository{db}
}

func (r *accessLogRepository) Create(ctx context.Context, log *domain.AccessLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}