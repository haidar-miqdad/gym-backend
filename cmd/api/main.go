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
	db.AutoMigrate(
		&domain.Member{},
		&domain.Package{},
		&domain.Subscription{}, 
		&domain.Payment{},
	)

	database.SeedPackages(db)

	
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- DEPENDENCY INJECTION ---
	// 1. Repository (Otot)
	memberRepo := repository.NewMemberRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// 2. Service (Otak)
	memberSvc := service.NewMemberService(memberRepo)
	// Kita berikan 'db' ke subSvc karena dia butuh mencari data Package secara langsung
	subSvc := service.NewSubscriptionService(subRepo, memberRepo, paymentRepo, db)

	// 3. Delivery (Mulut/API)
	delivery.NewMemberHandler(e, memberSvc)

	// 4. Initialize Handlers
	delivery.NewMemberHandler(e, memberSvc)
	delivery.NewSubscriptionHandler(e, subSvc)

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}