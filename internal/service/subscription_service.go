package service

import (
	"context"
	"errors"
	"fmt" // Digunakan untuk wrapping error agar lebih informatif
	"gym-backend/internal/domain"
	"gym-backend/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Interface harus sama persis dengan fungsi implementasi di bawah
type SubscriptionService interface {
	Subscribe(ctx context.Context, memberID, packageID, method, refNum string) (domain.Subscription, error)
}

type subscriptionService struct {
	subRepo    repository.SubscriptionRepository
	memberRepo repository.MemberRepository
	paymentRepo repository.PaymentRepository
	db         *gorm.DB
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository, memberRepo repository.MemberRepository, paymentRepo repository.PaymentRepository, db *gorm.DB) SubscriptionService {
	return &subscriptionService{
		subRepo:     subRepo,
		memberRepo:  memberRepo,
		paymentRepo: paymentRepo, // Masukkan ke struct
		db:          db,
	}
}

func (s *subscriptionService) Subscribe(ctx context.Context, memberID, packageID, method, refNum string) (domain.Subscription, error) {
	// 1. Validasi Format UUID
	mID, err := uuid.Parse(memberID)
	if err != nil {
		return domain.Subscription{}, errors.New("format ID member tidak valid")
	}

	pID, err := uuid.Parse(packageID)
	if err != nil {
		return domain.Subscription{}, errors.New("format ID paket tidak valid")
	}

	// 2. Cari Data Package
	var pkg domain.Package
	if err := s.db.WithContext(ctx).First(&pkg, "id = ?", pID).Error; err != nil {
		return domain.Subscription{}, errors.New("paket tidak ditemukan")
	}

	var newSub domain.Subscription

	// 3. Eksekusi Atomic Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		startDate := time.Now()
		endDate := startDate.AddDate(0, 0, pkg.DurationDays)

		// Logika khusus paket harian
		if pkg.DurationDays == 1 {
			endDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 23, 59, 59, 0, startDate.Location())
		}

		newSub = domain.Subscription{
			ID:        uuid.New(),
			MemberID:  mID,
			PackageID: pkg.ID,
			StartDate: startDate,
			EndDate:   endDate,
			Status:    "active",
		}

		if err := tx.Create(&newSub).Error; err != nil {
			return err
		}

		// Mencatat Payment otomatis saat subscribe
		newPayment := domain.Payment{
			ID:              uuid.New(),
			SubscriptionID:  newSub.ID,
			Amount:          pkg.Price,
			Method:          method,
			ReferenceNumber: refNum,
			Status:          "completed",
		}

		if err := s.paymentRepo.Create(ctx, tx, &newPayment); err != nil {
    return err
}

		return nil
	})

	if err != nil {
		return domain.Subscription{}, fmt.Errorf("gagal memproses transaksi: %w", err)
	}

	return newSub, nil
}