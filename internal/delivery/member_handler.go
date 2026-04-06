// Tugasnya: Menangani Request HTTP (Input/Output JSON).
package delivery

import (
	"gym-backend/internal/service"
	"net/http"
	"github.com/labstack/echo/v4"
)

// MemberHandler menangani request HTTP terkait entitas Member.
type MemberHandler struct {
	svc service.MemberService
}

func NewMemberHandler(public *echo.Group, protected *echo.Group, svc service.MemberService) {
	h := &MemberHandler{svc}
	
	// Rute ini tidak butuh token (Public Group)
	public.GET("/members/:id/status", h.GetStatus)
	
	// Rute ini wajib token (Protected Group)
	protected.POST("/members", h.RegisterMember)
	protected.GET("/members", h.GetAllMembers)
}

func (h *MemberHandler) RegisterMember(c echo.Context) error {
	// 1. Definisikan struktur request lokal (Data Transfer Object)
	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	// 2. Bind JSON dari request body ke struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Format request tidak valid",
		})
	}

	// 3. Panggil service layer
	member, err := h.svc.Register(c.Request().Context(), req.Name, req.Phone)
	if err != nil {
		// Mapping error bisnis ke status code HTTP yang sesuai
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}

	// 4. Return respon sukses (201 Created)
	return c.JSON(http.StatusCreated, member)
}

func (h *MemberHandler) GetAllMembers(c echo.Context) error {
	members, err := h.svc.GetAllMembers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil data member",
		})
	}

	return c.JSON(http.StatusOK, members)
}

// GET /api/v1/members/:id/status
func (h *MemberHandler) GetStatus(c echo.Context) error {
	id := c.Param("id")
	
	status, err := h.svc.GetMemberStatus(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, status)
}