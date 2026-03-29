package database

import (
	"gym-backend/internal/domain"
	"github.com/google/uuid"
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
			Name:         "Monthly Member",
			DurationDays: 30,
			Price:        300000,
			IsActive:     true,
		},
	}

	for _, pkg := range packages {
		// FirstOrCreate mencegah duplikasi data saat server restart
		db.Where(domain.Package{Name: pkg.Name}).FirstOrCreate(&pkg)
	}
}