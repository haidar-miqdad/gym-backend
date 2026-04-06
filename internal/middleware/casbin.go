package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/labstack/echo/v4"
)

func CheckPermission(enforcer *casbin.Enforcer, obj string, act string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Ambil roles (slice) dari context JWT
			rolesInterface := c.Get("roles")
			roles, ok := rolesInterface.([]interface{})
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Roles tidak ditemukan"})
			}

			// 2. Iterasi setiap role
			for _, r := range roles {
				roleName := r.(string)
				fmt.Printf("DEBUG: Mengecek Role: '%s' | Object: '%s' | Action: '%s'\n", roleName, obj, act)
				allowed, _ := enforcer.Enforce(roleName, obj, act)
				if allowed {
					// Jika salah satu role diizinkan, langsung lanjut ke handler
					return next(c)
				}
			}

			// 4. Jika semua role tidak ada yang cocok
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Akses ditolak: Tidak ada role Anda yang memiliki izin " + act,
			})
		}
	}
}