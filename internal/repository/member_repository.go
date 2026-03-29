// Tugasnya hanya satu: Berinteraksi langsung dengan Database SQL.
package repository

import (
	"context"
	"gym-backend/internal/domain"
	"gorm.io/gorm"
)

// MemberRepository mendefinisikan kontrak antara bisnis logic dan database.
// Kita menggunakan interface agar layer service tidak bergantung pada implementasi database tertentu (Decoupling).
type MemberRepository interface {
	Create(ctx context.Context, member *domain.Member) error
	FindByID(ctx context.Context, id string) (domain.Member, error)
	FindAll(ctx context.Context) ([]domain.Member, error)
	Update(ctx context.Context, member *domain.Member) error
	FindByPhone(ctx context.Context, phone string) (domain.Member, error)
}

type memberRepository struct {
	db *gorm.DB
}

// NewMemberRepository adalah Constructor untuk menginisialisasi repository.
func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{
		db: db,
	}
}

func (r *memberRepository) Create(ctx context.Context, member *domain.Member) error {
	// Menggunakan WithContext untuk mendukung cancellation dan tracing.
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *memberRepository) FindByID(ctx context.Context, id string) (domain.Member, error) {
	var member domain.Member
	err := r.db.WithContext(ctx).First(&member, "id = ?", id).Error
	return member, err
}

func (r *memberRepository) FindAll(ctx context.Context) ([]domain.Member, error) {
	var members []domain.Member
	err := r.db.WithContext(ctx).Find(&members).Error
	return members, err
}

func (r *memberRepository) Update(ctx context.Context, member *domain.Member) error {
	return r.db.WithContext(ctx).Save(member).Error
}


func (r *memberRepository) FindByPhone(ctx context.Context, phone string) (domain.Member, error) {
	var member domain.Member
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&member).Error
	return member, err
}