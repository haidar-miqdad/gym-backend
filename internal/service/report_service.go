package service

import (
	"context"
	"gym-backend/internal/repository"
	"time"
)

type ReportService interface {
	// Nama fungsi ini HARUS sesuai dengan yang dipanggil di handler
	GetDailyReport(ctx context.Context, date time.Time) (map[string]interface{}, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo}
}

func (s *reportService) GetDailyReport(ctx context.Context, date time.Time) (map[string]interface{}, error) {
	revenue, err := s.repo.GetDailyRevenue(ctx, date)
	if err != nil { return nil, err }

	attendance, err := s.repo.GetDailyAttendanceCount(ctx, date)
	if err != nil { return nil, err }

	return map[string]interface{}{
		"date":             date.Format("2006-01-02"),
		"daily_revenue":    revenue,
		"today_attendance": attendance,
	}, nil
}