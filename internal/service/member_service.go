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
	CheckIn(ctx context.Context, memberID string) (MemberStatusResponse, error)
	GetAllMembers(ctx context.Context, page, limit int) ([]domain.Member, error)
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
	if name == "" || phone == "" {
		return domain.Member{}, errors.New("nama dan nomor telepon wajib diisi")
	}

	_, err := s.repo.FindByPhone(ctx, phone)
	if err == nil {
		return domain.Member{}, fmt.Errorf("nomor telepon %s sudah terdaftar", phone)
	}
	
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Member{}, err
	}

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

// Tambahkan ctx ke parameter agar implementasi memenuhi interface
func (s *memberService) GetAllMembers(ctx context.Context, page, limit int) ([]domain.Member, error) {
	offset := (page - 1) * limit
	return s.repo.GetAll(ctx, limit, offset)
}

// Implementasi GetMemberStatus (Hanya Query/Read)
func (s *memberService) GetMemberStatus(ctx context.Context, memberID string) (MemberStatusResponse, error) {
	var sub domain.Subscription
	// Tambahkan Preload("Package") agar kita tahu nama paketnya
	err := s.db.WithContext(ctx).
		Preload("Package").
		Where("member_id = ? AND start_date <= ? AND end_date >= ?", memberID, time.Now(), time.Now()).
		First(&sub).Error

	if err != nil {
		return MemberStatusResponse{IsActive: false, Message: "Tidak ada paket aktif"}, nil
	}

	return MemberStatusResponse{
		IsActive:    true,
		Message:     "Paket Aktif",
		PackageName: sub.Package.Name,
		EndDate:     sub.EndDate,
		DaysLeft:    int(time.Until(sub.EndDate).Hours() / 24),
	}, nil
}

// Implementasi CheckIn (Command/Write)
func (s *memberService) CheckIn(ctx context.Context, memberID string) (MemberStatusResponse, error) {
	// 1. Cek status (Gunakan fungsi yang sudah ada)
	status, err := s.GetMemberStatus(ctx, memberID)
	if err != nil || !status.IsActive {
		return status, errors.New("akses ditolak: " + status.Message)
	}

	// 2. Ambil data subscription ID untuk logging
	var sub domain.Subscription
	s.db.Where("member_id = ? AND status = 'active'", memberID).First(&sub)

	// 3. Catat ke AccessLog (Hanya dilakukan di fungsi POST ini)
	mID, _ := uuid.Parse(memberID)
	log := domain.AccessLog{
		ID:             uuid.New(),
		MemberID:       mID,
		SubscriptionID: sub.ID,
		CheckInAt:      time.Now(),
	}
	
	if err := s.accessLogRepo.Create(ctx, &log); err != nil {
		return status, fmt.Errorf("gagal mencatat kehadiran: %w", err)
	}

	return status, nil
}