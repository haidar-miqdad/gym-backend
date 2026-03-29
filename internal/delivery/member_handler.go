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

// NewMemberHandler menginisialisasi rute-rute API untuk Member.
func NewMemberHandler(e *echo.Echo, svc service.MemberService) {
	h := &MemberHandler{
		svc: svc,
	}

	// Grouping API v1
	api := e.Group("/api/v1")
	api.POST("/members", h.RegisterMember)
	api.GET("/members", h.GetAllMembers)
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