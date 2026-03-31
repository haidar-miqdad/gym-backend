// Tugasnya: Logika Bisnis (misal: validasi jika nomor telepon sudah terdaftar).
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

type MemberService interface {
	Register(ctx context.Context, name, phone string) (domain.Member, error)
	GetAllMembers(ctx context.Context) ([]domain.Member, error)
	GetMemberStatus(ctx context.Context, memberID string) (MemberStatusResponse, error)
}

type memberService struct {
	repo repository.MemberRepository
	accessLogRepo repository.AccessLogRepository
	db   *gorm.DB
}

type MemberStatusResponse struct {
	IsActive    bool      `json:"is_active"`
	Message     string    `json:"message"`
	PackageName string    `json:"package_name,omitempty"`
	DaysLeft    int       `json:"days_left"`
	EndDate     time.Time `json:"end_date"`
}

func NewMemberService(repo repository.MemberRepository, logRepo repository.AccessLogRepository, db *gorm.DB) MemberService {
	return &memberService{
		repo:          repo,
		accessLogRepo: logRepo,
		db:            db,
	}
}

func (s *memberService) Register(ctx context.Context, name, phone string) (domain.Member, error) {
	// 1. Validasi Input (WAJIB pakai IF dan RETURN)
	if name == "" || phone == "" {
		return domain.Member{}, errors.New("nama dan nomor telepon wajib diisi")
	}

	// 2. Cek Duplikasi (Pakai "_" karena existingMember tidak kita pakai datanya)
	_, err := s.repo.FindByPhone(ctx, phone)
	
	// Jika err == nil, artinya data DITEMUKAN (berarti duplikat)
	if err == nil {
		return domain.Member{}, fmt.Errorf("nomor telepon %s sudah terdaftar", phone)
	}
	
	// Jika error-nya bukan 'Record Not Found', berarti ada masalah di database
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Member{}, err
	}

	// 3. Buat Member Baru
	newMember := domain.Member{
		ID:     uuid.New(),
		Name:   name,
		Phone:  phone,
		Status: "active",
	}

	if err := s.repo.Create(ctx, &newMember); err != nil {
		return domain.Member{}, fmt.Errorf("gagal simpan ke database: %w", err)
	}

	return newMember, nil
}

func (s *memberService) GetAllMembers(ctx context.Context) ([]domain.Member, error) {
	return s.repo.FindAll(ctx)
}

func (s *memberService) GetMemberStatus(ctx context.Context, memberID string) (MemberStatusResponse, error) {
	// 1. Cek apakah member ada
	_, err := s.repo.GetByID(ctx, memberID)
	if err != nil {
		return MemberStatusResponse{}, errors.New("member tidak ditemukan")
	}

	// 2. Cari subscription aktif melalui subRepo
	// (Catatan: Pastikan memberService punya akses ke subRepo atau gunakan DB langsung)
	var sub domain.Subscription
	err = s.db.Where("member_id = ? AND start_date <= ? AND end_date >= ?", 
		memberID, time.Now(), time.Now()).First(&sub).Error

	if err != nil {
		return MemberStatusResponse{
			IsActive: false,
			Message:  "Akses Ditolak: Tidak ada paket aktif atau sudah expired",
		}, nil
	}

	// 3. Hitung sisa hari
	daysLeft := int(time.Until(sub.EndDate).Hours() / 24)
	if daysLeft < 0 { daysLeft = 0 }

	// 4. Buat response
	mID, err := uuid.Parse(memberID)
	if err != nil {
		return MemberStatusResponse{}, errors.New("format ID member tidak valid")
	}

	go func() {
        log := domain.AccessLog{
            ID:             uuid.New(),
            MemberID:       mID,
            SubscriptionID: sub.ID,
            CheckInAt:      time.Now(),
        }
        s.accessLogRepo.Create(context.Background(), &log)
    }()

	return MemberStatusResponse{
		IsActive:    true,
		Message:     "Akses Diterima",
		DaysLeft:    daysLeft,
		EndDate:     sub.EndDate,
	}, nil
}