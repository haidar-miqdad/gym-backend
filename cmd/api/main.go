package main

import (
	"gym-backend/internal/delivery"
	"gym-backend/internal/domain"
	"gym-backend/internal/repository"
	"gym-backend/internal/service"
	"gym-backend/pkg/database"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	db := database.InitDB()
	db.AutoMigrate(&domain.Member{}, &domain.Package{}, &domain.Subscription{})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- DEPENDENCY INJECTION ---
	// 1. Repository (Otot)
	memberRepo := repository.NewMemberRepository(db)

	// 2. Service (Otak)
	memberSvc := service.NewMemberService(memberRepo)

	// 3. Delivery (Mulut/API)
	delivery.NewMemberHandler(e, memberSvc)

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}