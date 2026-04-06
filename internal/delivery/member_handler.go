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
    
    public.GET("/members/:id/status", h.GetStatus)
    protected.POST("/members", h.RegisterMember, middleware.CheckPermission(enforcer, "members", "create"))
    protected.GET("/members", h.GetAllMembers, middleware.CheckPermission(enforcer, "members", "view"))
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

	// Tambahkan c.Request().Context() sebagai argumen pertama
	members, err := h.svc.GetAllMembers(c.Request().Context(), page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) GetStatus(c echo.Context) error {
	id := c.Param("id")
	status, err := h.svc.GetMemberStatus(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, status)
}

