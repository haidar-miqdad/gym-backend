package main

import (
	"gym-backend/internal/delivery"
	"gym-backend/internal/domain"
	"gym-backend/internal/middleware"
	"gym-backend/internal/repository"
	"gym-backend/internal/service"
	"gym-backend/pkg/database"
	"os"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)


func main() {
	godotenv.Load()
	db := database.InitDB()
	
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

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
	database.SeedStaff(db)

	// 2. Inisialisasi Casbin
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		e.Logger.Fatal("Gagal membuat adapter Casbin:", err)
	}

	enforcer, err := casbin.NewEnforcer("rbac_model.conf", adapter)
	if err != nil {
		e.Logger.Fatal("Gagal membuat enforcer Casbin:", err)
	}
	enforcer.LoadPolicy()

	// --- BAGIAN PERUBAHAN: IDEMPOTENT POLICY SEEDING ---

	// Izin Mengelola Permission (Hanya Admin)
	hasManagePermission, _ := enforcer.HasPolicy("admin", "permissions", "manage")
	if !hasManagePermission {
		enforcer.AddPolicy("admin", "permissions", "manage")
	}

	// Izin Laporan (Hanya Admin)
	hasReportPolicy, _ := enforcer.HasPolicy("admin", "reports", "view")
	if !hasReportPolicy {
		enforcer.AddPolicy("admin", "reports", "view")
	}

	// Izin Member (Hanya Staff)
	hasStaffPolicy, _ := enforcer.HasPolicy("staff", "members", "view")
	if !hasStaffPolicy {
		enforcer.AddPolicy("staff", "members", "view")
	}

	// Izin Membuat Member (Admin)
	hasMemberPolicy, _ := enforcer.HasPolicy("admin", "members", "create")
	if !hasMemberPolicy {
		enforcer.AddPolicy("admin", "members", "create")
	}

	hasSubPolicy, _ := enforcer.HasPolicy("admin", "subscriptions", "create")
	if !hasSubPolicy {
		enforcer.AddPolicy("admin", "subscriptions", "create")
	}

	hasCheckInPolicy, _ := enforcer.HasPolicy("admin", "attendance", "create")
	if !hasCheckInPolicy {
		enforcer.AddPolicy("admin", "attendance", "create")
	}

	hasViewReport, _ := enforcer.HasPolicy("admin", "reports", "view")
	if !hasViewReport {
		enforcer.AddPolicy("admin", "reports", "view")
	}

	// --- PERUBAHAN KRUSIAL: RBAC INHERITANCE (GROUPING POLICY) ---
	// Membuat Admin otomatis memiliki semua izin yang dimiliki Staff
	hasAdminStaffLink, _ := enforcer.HasGroupingPolicy("admin", "staff")
	if !hasAdminStaffLink {
		enforcer.AddGroupingPolicy("admin", "staff")
	}

	// 4. Dependency Injection (DI)
	memberRepo := repository.NewMemberRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	accessLogRepo := repository.NewAccessLogRepository(db)
	reportRepo := repository.NewReportRepository(db)

	memberSvc := service.NewMemberService(memberRepo, accessLogRepo, db)
	subSvc := service.NewSubscriptionService(subRepo, memberRepo, paymentRepo, db)
	reportSvc := service.NewReportService(reportRepo)
	authSvc := service.NewAuthService(db)

	// 5. Route Grouping Logic
	v1 := e.Group("/api/v1")
	protected := v1.Group("") 
	protected.Use(middleware.JWTMiddleware)

	// 6. Initialize Handlers
	delivery.NewAuthHandler(v1, authSvc)
	delivery.NewReportHandler(protected, reportSvc, enforcer) 
	delivery.NewSubscriptionHandler(protected, subSvc, enforcer)
	delivery.NewMemberHandler(v1, protected, memberSvc, enforcer)
	delivery.NewPermissionHandler(protected, enforcer)
	delivery.NewAccessLogHandler(protected, memberSvc)

	port := os.Getenv("APP_PORT")
	if port == "" { port = "8080" }
	e.Logger.Fatal(e.Start(":" + port))
}