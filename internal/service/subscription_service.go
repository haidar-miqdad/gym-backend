package service

import (
	"context"
	"errors"
	"fmt"
	"gym-backend/internal/domain"
	"gym-backend/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionService interface {
	Subscribe(ctx context.Context, memberID, packageID string) (domain.Subscription, error)
}

type subscriptionService struct {
	subRepo    repository.SubscriptionRepository
	memberRepo repository.MemberRepository
	// Kita butuh PackageRepo (asumsi sudah buat atau gunakan GORM langsung untuk ringkasnya)
	db *gorm.DB 
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository, memberRepo repository.MemberRepository, db *gorm.DB) SubscriptionService {
	return &subscriptionService{subRepo, memberRepo, db}
}

func (s *subscriptionService) Subscribe(ctx context.Context, memberID, packageID string) (domain.Subscription, error) {
	// 1. Validasi Format UUID untuk MemberID
	mID, err := uuid.Parse(memberID)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("format ID member tidak valid: %w", err)
	}

	// 2. Cari Data Package (Cek harga & durasi)
	var pkg domain.Package
	if err := s.db.WithContext(ctx).First(&pkg, "id = ?", packageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Subscription{}, errors.New("paket tidak ditemukan")
		}
		return domain.Subscription{}, err
	}

	// 3. Hitung Tanggal (Business Logic)
	startDate := time.Now()
	var endDate time.Time

	if pkg.DurationDays == 1 {
		// Kasus Harian: Expired jam 23:59:59 hari ini
		endDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 23, 59, 59, 0, startDate.Location())
	} else {
		// Kasus Bulanan/Tahunan: Tambah hari secara normal
		endDate = startDate.AddDate(0, 0, pkg.DurationDays)
	}

	// 4. Mapping ke Model (Sekarang menggunakan mID yang sudah divalidasi)
	newSub := domain.Subscription{
		ID:        uuid.New(),
		MemberID:  mID, // Menggunakan hasil parse yang sukses
		PackageID: pkg.ID,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "active",
	}

	// 5. Simpan (Gunakan Repository)
	if err := s.subRepo.Create(ctx, &newSub); err != nil {
		return domain.Subscription{}, fmt.Errorf("gagal memproses langganan: %w", err)
	}

	return newSub, nil
}