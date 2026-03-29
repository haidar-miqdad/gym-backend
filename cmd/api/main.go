package main

import (
	"gym-backend/internal/domain" // Pastikan import domain masuk
	"gym-backend/pkg/database"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load environment variables
	godotenv.Load()

	// 2. Initialize Database
	db := database.InitDB()
    
	// 3. AKTIFKAN INI: Jalankan Migrasi
	// Ini akan menggunakan variabel 'db' sehingga error "unused" hilang
	db.AutoMigrate(&domain.Member{})

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Basic Health Check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "UP",
			"message": "Gym System Backend is Running",
		})
	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}