package delivery

import (
	"gym-backend/internal/middleware"
	"gym-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/casbin/casbin/v3"
	"github.com/labstack/echo/v4"
)

type MemberHandler struct {
	svc service.MemberService
}

func NewMemberHandler(public *echo.Group, protected *echo.Group, svc service.MemberService, enforcer *casbin.Enforcer) {
    h := &MemberHandler{svc}
    
    // GET: Hanya untuk melihat status (misal: di aplikasi member atau dashboard)
    public.GET("/members/:id/status", h.GetStatus)
    
    // POST: Untuk mesin absen / gate gym
    protected.POST("/check-in", h.CheckIn, middleware.CheckPermission(enforcer, "attendance", "create"))
}

func (h *MemberHandler) CheckIn(c echo.Context) error {
	var req struct {
		MemberID string `json:"member_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request salah"})
	}

	status, err := h.svc.CheckIn(c.Request().Context(), req.MemberID)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, status)
}

func (h *MemberHandler) RegisterMember(c echo.Context) error {
	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
	}

	member, err := h.svc.Register(c.Request().Context(), req.Name, req.Phone)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, member)
}

func (h *MemberHandler) GetAllMembers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 { limit = 10 }

	members, err := h.svc.GetAllMembers(c.Request().Context(), page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Response dengan Metadata untuk memudahkan Flutter/React
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": members,
		"meta": map[string]interface{}{
			"current_page": page,
			"limit":        limit,
			"total_count":  len(members), // Opsional: Tambahkan Count query di repo untuk hasil akurat
		},
	})
}

func (h *MemberHandler) GetStatus(c echo.Context) error {
	id := c.Param("id")
	status, err := h.svc.GetMemberStatus(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, status)
}

