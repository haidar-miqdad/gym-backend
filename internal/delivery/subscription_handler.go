package delivery

import (
	"gym-backend/internal/service"
	"net/http"
	"github.com/labstack/echo/v4"
)

type SubscriptionHandler struct {
	svc service.SubscriptionService
}

func NewSubscriptionHandler(e *echo.Echo, svc service.SubscriptionService) {
	h := &SubscriptionHandler{svc}
	
	api := e.Group("/api/v1")
	api.POST("/subscriptions", h.Subscribe)
}

func (h *SubscriptionHandler) Subscribe(c echo.Context) error {
	var req struct {
		MemberID  string `json:"member_id"`
		PackageID string `json:"package_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	sub, err := h.svc.Subscribe(c.Request().Context(), req.MemberID, req.PackageID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, sub)
}