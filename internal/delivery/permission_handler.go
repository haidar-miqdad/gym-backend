package delivery

import (
	"net/http"
	"github.com/casbin/casbin/v3"
	"github.com/labstack/echo/v4"
	"gym-backend/internal/middleware"
)

type PermissionHandler struct {
	enforcer *casbin.Enforcer
}

func NewPermissionHandler(g *echo.Group, enforcer *casbin.Enforcer) {
	h := &PermissionHandler{enforcer}
	
	// Hanya Super Admin yang boleh mengakses endpoint ini
	permissionGroup := g.Group("/permissions")
	permissionGroup.Use(middleware.CheckPermission(enforcer, "permissions", "manage"))

	permissionGroup.POST("", h.AddPermission)
	permissionGroup.DELETE("", h.RemovePermission)
	permissionGroup.GET("/:role", h.GetPermissionsByRole)
}

type PermissionRequest struct {
	Role     string `json:"role"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

func (h *PermissionHandler) AddPermission(c echo.Context) error {
	var req PermissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request salah"})
	}

	// Menambahkan aturan ke database
	added, err := h.enforcer.AddPolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !added {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Izin sudah ada"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Izin berhasil ditambahkan"})
}

func (h *PermissionHandler) RemovePermission(c echo.Context) error {
	var req PermissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format request salah"})
	}

	removed, _ := h.enforcer.RemovePolicy(req.Role, req.Resource, req.Action)
	if !removed {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Izin tidak ditemukan"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Izin berhasil dihapus"})
}

func (h *PermissionHandler) GetPermissionsByRole(c echo.Context) error {
	role := c.Param("role")
	
	// Casbin v3 mengembalikan ([][]string, error)
	policies, err := h.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil data kebijakan: " + err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, policies)
}