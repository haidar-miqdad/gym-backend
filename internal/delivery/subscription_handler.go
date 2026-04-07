package delivery

import (
	"gym-backend/internal/middleware"
	"gym-backend/internal/service"
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/labstack/echo/v4"
)

// Subscribe godoc
// @Summary Mendaftarkan member ke paket gym
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body SubscribeRequest true "Data Langganan"
// @Success 201 {object} domain.Subscription
// @Failure 403 {object} map[string]string
// @Router /subscriptions [post]
// @Security ApiKeyAuth

// SubscriptionHandler menangani request HTTP untuk transaksi langganan.
type SubscriptionHandler struct {
	svc service.SubscriptionService
}

// NewSubscriptionHandler menginisialisasi rute untuk subscription.
func NewSubscriptionHandler(protected *echo.Group, svc service.SubscriptionService, enforcer *casbin.Enforcer) {
	h := &SubscriptionHandler{svc}

	// Grup rute subscription dengan proteksi Casbin
	// Izin: Role harus memiliki permission 'create' pada resource 'subscriptions'
	subGroup := protected.Group("/subscriptions")
	subGroup.POST("", h.Subscribe, middleware.CheckPermission(enforcer, "subscriptions", "create"))
}

// SubscribeRequest mendefinisikan struktur input JSON dari klien.
type SubscribeRequest struct {
	MemberID        string `json:"member_id"`
	PackageID       string `json:"package_id"`
	Method          string `json:"method"`           // misal: cash, transfer, midtrans
	ReferenceNumber string `json:"reference_number"` // kode unik transaksi
}

func (h *SubscriptionHandler) Subscribe(c echo.Context) error {
	var req SubscribeRequest

	// 1. Bind JSON request body ke struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Format request tidak valid",
		})
	}

	// 2. Validasi input sederhana
	if req.MemberID == "" || req.PackageID == "" || req.Method == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Member ID, Package ID, dan Method wajib diisi",
		})
	}

	// 3. Panggil service layer untuk memproses langganan & pembayaran
	subscription, err := h.svc.Subscribe(
		c.Request().Context(),
		req.MemberID,
		req.PackageID,
		req.Method,
		req.ReferenceNumber,
	)

	if err != nil {
		// Mengembalikan error bisnis (misal: paket tidak ditemukan/format UUID salah)
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}

	// 4. Return respon sukses (201 Created)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Langganan berhasil didaftarkan",
		"data":    subscription,
	})
}