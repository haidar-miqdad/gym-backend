package service

import (
	"context"
	"gym-backend/internal/repository"
	"time"
)

type DailyReportResponse struct {
	Date            string  `json:"date"`
	TotalRevenue    float64 `json:"total_revenue"`
	TotalAttendance int64   `json:"total_attendance"`
}

type ReportService interface {
	GetDailyReport(ctx context.Context, date time.Time) (DailyReportResponse, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo}
}

func (s *reportService) GetDailyReport(ctx context.Context, date time.Time) (DailyReportResponse, error) {
	revenue, err := s.repo.GetTotalRevenueByDate(ctx, date)
	if err != nil {
		return DailyReportResponse{}, err
	}

	attendance, err := s.repo.GetTotalAttendanceByDate(ctx, date)
	if err != nil {
		return DailyReportResponse{}, err
	}

	return DailyReportResponse{
		Date:            date.Format("2006-01-02"),
		TotalRevenue:    revenue,
		TotalAttendance: attendance,
	}, nil
}