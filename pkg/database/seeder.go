package database

import (
	"gym-backend/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedPackages(db *gorm.DB) {
	packages := []domain.Package{
		{
			ID:           uuid.New(),
			Name:         "Daily Pass",
			DurationDays: 1,
			Price:        35000,
			IsActive:     true,
		},
		{
			ID:           uuid.New(),
			Name:         "Monthly Membership",
			DurationDays: 30,
			Price:        300000,
			IsActive:     true,
		},
	}

	for _, pkg := range packages {
		// Cek berdasarkan Nama, jika belum ada baru buat
		db.Where(domain.Package{Name: pkg.Name}).FirstOrCreate(&pkg)
	}
}

func SeedAdmin(db *gorm.DB) {
    hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    admin := domain.User{
        Username: "admin",
        Password: string(hash),
        Role:     "admin",
    }

    // Gunakan ini untuk mencegah error duplicate key
    db.Where(domain.User{Username: "admin"}).FirstOrCreate(&admin)
}

func SeedStaff(db *gorm.DB) {
    hash, _ := bcrypt.GenerateFromPassword([]byte("staff123"), bcrypt.DefaultCost)
    staff := domain.User{
        Username: "staff",
        Password: string(hash),
        Role:     "staff",
    }
    // Gunakan FirstOrCreate untuk menghindari error duplicate key yang Anda alami sebelumnya
    db.Where(domain.User{Username: "staff"}).FirstOrCreate(&staff)
}