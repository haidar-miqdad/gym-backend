package repository

import (
	"context"
	"time"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetTotalRevenueByDate(ctx context.Context, date time.Time) (float64, error)
	GetTotalAttendanceByDate(ctx context.Context, date time.Time) (int64, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) GetTotalRevenueByDate(ctx context.Context, date time.Time) (float64, error) {
	var total float64
	// Query: SELECT SUM(amount) FROM payments WHERE created_at::date = ?
	err := r.db.WithContext(ctx).Table("payments").
		Where("DATE(created_at) = ?", date.Format("2006-01-02")).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *reportRepository) GetTotalAttendanceByDate(ctx context.Context, date time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("access_logs").
		Where("DATE(check_in_at) = ?", date.Format("2006-01-02")).
		Count(&count).Error
	return count, err
}