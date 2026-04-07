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

// SubscriptionService mendefinisikan kontrak utama untuk manajemen langganan.
type SubscriptionService interface {
	Subscribe(ctx context.Context, memberID, packageID, method, refNum string) (domain.Subscription, error)
}

type subscriptionService struct {
	subRepo     repository.SubscriptionRepository
	memberRepo  repository.MemberRepository
	paymentRepo repository.PaymentRepository
	db          *gorm.DB
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository, memberRepo repository.MemberRepository, paymentRepo repository.PaymentRepository, db *gorm.DB) SubscriptionService {
	return &subscriptionService{
		subRepo:     subRepo,
		memberRepo:  memberRepo,
		paymentRepo: paymentRepo,
		db:          db,
	}
}

func (s *subscriptionService) Subscribe(ctx context.Context, memberID, packageID, method, refNum string) (domain.Subscription, error) {
	// 1. Validasi Format UUID untuk Member dan Package
	mID, err := uuid.Parse(memberID)
	if err != nil {
		return domain.Subscription{}, errors.New("format ID member tidak valid")
	}

	pID, err := uuid.Parse(packageID)
	if err != nil {
		return domain.Subscription{}, errors.New("format ID paket tidak valid")
	}

	// 2. Cari Data Package untuk mendapatkan harga dan durasi
	var pkg domain.Package
	if err := s.db.WithContext(ctx).First(&pkg, "id = ?", pID).Error; err != nil {
		return domain.Subscription{}, errors.New("paket tidak ditemukan")
	}

	var newSub domain.Subscription
	var lastSub domain.Subscription
	var startDate time.Time

	// 3. Logika Stacking: Cek apakah ada langganan yang masih aktif
	err = s.db.WithContext(ctx).
		Where("member_id = ? AND end_date > ?", mID, time.Now()).
		Order("end_date DESC").
		First(&lastSub).Error

	if err == nil {
		// Jika ada, langganan baru mulai tepat setelah langganan lama habis
		startDate = lastSub.EndDate
	} else {
		// Jika tidak ada atau sudah expired, mulai dari detik ini
		startDate = time.Now()
	}

	// 4. Hitung Tanggal Berakhir (Handle durasi harian vs bulanan)
	endDate := startDate.AddDate(0, 0, pkg.DurationDays)
	
	if pkg.DurationDays == 1 {
		// Khusus paket 1 hari, expired di akhir hari tersebut
		endDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 23, 59, 59, 0, startDate.Location())
	}

	// 5. Eksekusi Atomic Transaction (All-or-Nothing)
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// A. Buat Record Subscription
		newSub = domain.Subscription{
			ID:        uuid.New(),
			MemberID:  mID,
			PackageID: pkg.ID,
			StartDate: startDate,
			EndDate:   endDate,
			Status:    "active",
		}

		if err := tx.Create(&newSub).Error; err != nil {
			return fmt.Errorf("gagal membuat subscription: %w", err)
		}

		// B. Buat Record Payment
		newPayment := domain.Payment{
			ID:              uuid.New(),
			SubscriptionID:  newSub.ID,
			Amount:          pkg.Price,
			Method:          method,
			ReferenceNumber: refNum,
			Status:          "completed",
		}

		// Pastikan repository payment Anda menerima parameter 'tx' untuk menjaga atomisitas
		if err := s.paymentRepo.Create(ctx, tx, &newPayment); err != nil {
			return fmt.Errorf("gagal membuat payment: %w", err)
		}

		// C. Update status Member menjadi 'active' secara otomatis
		if err := tx.Model(&domain.Member{}).Where("id = ?", mID).Update("status", "active").Error; err != nil {
			return fmt.Errorf("gagal memperbarui status member: %w", err)
		}

		return nil
	})

	if err != nil {
		return domain.Subscription{}, fmt.Errorf("transaksi gagal: %w", err)
	}

	return newSub, nil
}