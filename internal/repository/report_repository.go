package repository

import (
	"context"
	"time"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetDailyRevenue(ctx context.Context, date time.Time) (float64, error)
	GetDailyAttendanceCount(ctx context.Context, date time.Time) (int64, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) GetDailyRevenue(ctx context.Context, date time.Time) (float64, error) {
	var total float64
	// Filter berdasarkan hari yang dipilih menggunakan DATE() di PostgreSQL
	err := r.db.WithContext(ctx).Table("payments").
		Where("status = ? AND DATE(created_at) = DATE(?)", "completed", date).
		Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *reportRepository) GetDailyAttendanceCount(ctx context.Context, date time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("access_logs").
		Where("DATE(check_in_at) = DATE(?)", date).
		Count(&count).Error
	return count, err
}