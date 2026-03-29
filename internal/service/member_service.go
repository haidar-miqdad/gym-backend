// Tugasnya: Logika Bisnis (misal: validasi jika nomor telepon sudah terdaftar).
package service

import (
	"context"
	"errors"
	"fmt"
	"gym-backend/internal/domain"
	"gym-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberService interface {
	Register(ctx context.Context, name, phone string) (domain.Member, error)
	GetAllMembers(ctx context.Context) ([]domain.Member, error)
}

type memberService struct {
	repo repository.MemberRepository
}

func NewMemberService(repo repository.MemberRepository) MemberService {
	return &memberService{
		repo: repo,
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