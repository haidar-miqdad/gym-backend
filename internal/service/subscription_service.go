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

	// 3. LOGIKA STACKING (Harus SEBELUM Transaction)
	var lastSub domain.Subscription
	var startDate time.Time

	// Cari apakah ada langganan yang masih aktif
	err = s.db.WithContext(ctx).
		Where("member_id = ? AND end_date > ?", mID, time.Now()).
		Order("end_date DESC").
		First(&lastSub).Error

	if err == nil {
		// Jika ada, langganan baru mulai saat langganan lama habis
		startDate = lastSub.EndDate
	} else {
		// Jika tidak ada/expired, mulai dari sekarang
		startDate = time.Now()
	}

	// 4. Hitung EndDate berdasarkan durasi paket
	endDate := startDate.AddDate(0, 0, pkg.DurationDays)
	
	// Logic khusus paket harian (expired jam 23:59 hari yang sama)
	// Catatan: Jika harian di-stack, dia akan expired di 23:59 di hari startDate tersebut
	if pkg.DurationDays == 1 {
		endDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 23, 59, 59, 0, startDate.Location())
	}

	// 5. EKSEKUSI ATOMIC TRANSACTION
	err = s.db.Transaction(func(tx *gorm.DB) error {
		newSub = domain.Subscription{
			ID:        uuid.New(),
			MemberID:  mID,
			PackageID: pkg.ID,
			StartDate: startDate, // Menggunakan startDate hasil stacking
			EndDate:   endDate,   // Menggunakan endDate hasil stacking
			Status:    "active",
		}

		if err := tx.Create(&newSub).Error; err != nil {
			return err
		}

		newPayment := domain.Payment{
			ID:              uuid.New(),
			SubscriptionID:  newSub.ID,
			Amount:          pkg.Price,
			Method:          method,
			ReferenceNumber: refNum,
			Status:          "completed",
		}

		// Menggunakan repository untuk create payment
		return s.paymentRepo.Create(ctx, tx, &newPayment)
	})

	if err != nil {
		return domain.Subscription{}, fmt.Errorf("gagal memproses transaksi: %w", err)
	}

	return newSub, nil
}