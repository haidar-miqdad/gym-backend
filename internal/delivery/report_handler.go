package delivery

import (
	"gym-backend/internal/service"
	"net/http"
	"time"
	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	svc service.ReportService
}

func NewReportHandler(e *echo.Echo, svc service.ReportService) {
	h := &ReportHandler{svc}
	
	api := e.Group("/api/v1/reports")
	api.GET("/daily", h.GetDailyReport)
}

func (h *ReportHandler) GetDailyReport(c echo.Context) error {
	// Default ke hari ini jika tidak ada parameter tanggal
	dateStr := c.QueryParam("date")
	var date time.Time
	var err error

	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format tanggal harus YYYY-MM-DD"})
		}
	}

	report, err := h.svc.GetDailyReport(c.Request().Context(), date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, report)
}