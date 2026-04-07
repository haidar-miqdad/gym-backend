package delivery

import (
	"gym-backend/internal/service"
	"net/http"
	"github.com/labstack/echo/v4"
)

type AccessLogHandler struct {
	svc service.MemberService // Sementara menggunakan member service atau service khusus log
}

func NewAccessLogHandler(protected *echo.Group, svc service.MemberService) {
	h := &AccessLogHandler{svc}
	protected.GET("/access-logs", h.GetLogs)
}

func (h *AccessLogHandler) GetLogs(c echo.Context) error {
	// Implementasi query log akses Anda di sini
	return c.JSON(http.StatusOK, map[string]string{"message": "Endpoint log akses siap"})
}