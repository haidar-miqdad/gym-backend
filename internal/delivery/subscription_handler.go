package delivery

import (
	"gym-backend/internal/middleware"
	"gym-backend/internal/service"
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/labstack/echo/v4"
)

type SubscriptionHandler struct {
	svc service.SubscriptionService
}

func NewSubscriptionHandler(g *echo.Group, svc service.SubscriptionService, enforcer *casbin.Enforcer) {
	h := &SubscriptionHandler{svc}
	
	g.POST("/subscriptions", h.Subscribe, middleware.CheckPermission(enforcer, "subscriptions", "create"))
}

func (h *SubscriptionHandler) Subscribe(c echo.Context) error {
	// Definisi struct request untuk menangkap input JSON
	var req struct {
		MemberID  string `json:"member_id"`
		PackageID string `json:"package_id"`
		Method    string `json:"method"`
		Reference string `json:"reference_number"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	// Memanggil service dengan 5 parameter sesuai kontrak baru
	sub, err := h.svc.Subscribe(
		c.Request().Context(), 
		req.MemberID, 
		req.PackageID, 
		req.Method, 
		req.Reference,
	)
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, sub)
}