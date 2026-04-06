package main

import (
	"gym-backend/internal/delivery"
	"gym-backend/internal/domain"
	"gym-backend/internal/middleware"
	"gym-backend/internal/repository"
	"gym-backend/internal/service"
	"gym-backend/pkg/database"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	db := database.InitDB()
	
	// 1. Database Migration & Seeding
	db.AutoMigrate(
		&domain.Member{},
		&domain.Package{},
		&domain.Subscription{}, 
		&domain.Payment{},
		&domain.AccessLog{},
		&domain.User{},
	)
	database.SeedPackages(db)
	database.SeedAdmin(db)

	// 2. Echo Instance & Global Middleware
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// 3. Dependency Injection (DI)
	// Repositories
	memberRepo := repository.NewMemberRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	accessLogRepo := repository.NewAccessLogRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Services
	memberSvc := service.NewMemberService(memberRepo, accessLogRepo, db)
	subSvc := service.NewSubscriptionService(subRepo, memberRepo, paymentRepo, db)
	reportSvc := service.NewReportService(reportRepo)
	authSvc := service.NewAuthService(db)

	// 4. Route Grouping Logic (Nesting Strategy)
// Kita buat base group v1 sebagai "payung" besar
v1 := e.Group("/api/v1")

// Kita buat sub-group "protected" di bawah v1. 
// Path kosong "" berarti prefixnya tetap /api/v1
protected := v1.Group("") 
protected.Use(middleware.JWTMiddleware)

// 5. Initialize Handlers
// Handler Login: Daftarkan ke v1 (Public / Tanpa Gembok)
delivery.NewAuthHandler(v1, authSvc)

// Handler Report & Subscription: Daftarkan ke protected (Wajib Token)
delivery.NewReportHandler(protected, reportSvc)
delivery.NewSubscriptionHandler(protected, subSvc)

// Handler Member: v1 untuk rute publik (status), protected untuk rute admin (CRUD)
delivery.NewMemberHandler(v1, protected, memberSvc)

	port := os.Getenv("APP_PORT")
if port == "" { port = "8080" }
e.Logger.Fatal(e.Start(":" + port))
}