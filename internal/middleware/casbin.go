package middleware

import (
	"net/http"
	"github.com/casbin/casbin/v3"
	"github.com/labstack/echo/v4"
)

func CheckPermission(enforcer *casbin.Enforcer, obj string, act string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Ambil role yang disuntikkan JWTMiddleware sebelumnya
			role, ok := c.Get("role").(string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Role tidak ditemukan"})
			}

			// 2. Tanya ke Casbin: "Apakah role [admin] boleh [view] pada [reports]?"
			// Casbin akan mencocokkan dengan data di TablePlus Anda tadi
			ok, err := enforcer.Enforce(role, obj, act)
			if err != nil || !ok {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Akses ditolak: Anda tidak punya izin " + act + " untuk " + obj,
				})
			}

			return next(c)
		}
	}
}