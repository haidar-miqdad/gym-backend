package repository

import (
	"context"
	"gym-backend/internal/domain"
	"gorm.io/gorm"
)

type MemberRepository interface {
	// Tambahkan ctx agar sinkron dengan metode lain
	GetAll(ctx context.Context, limit, offset int) ([]domain.Member, error)
	Create(ctx context.Context, member *domain.Member) error
	FindByPhone(ctx context.Context, phone string) (domain.Member, error)
	GetByID(ctx context.Context, id string) (domain.Member, error)
	// Hapus FindByID/FindAll yang duplikat jika tidak digunakan untuk kebersihan kode
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Member, error) {
	var members []domain.Member
	// Tambahkan WithContext(ctx)
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&members).Error
	return members, err
}

func (r *memberRepository) Create(ctx context.Context, member *domain.Member) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *memberRepository) FindByPhone(ctx context.Context, phone string) (domain.Member, error) {
	var member domain.Member
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&member).Error
	return member, err
}

func (r *memberRepository) GetByID(ctx context.Context, id string) (domain.Member, error) {
    var member domain.Member
    err := r.db.WithContext(ctx).First(&member, "id = ?", id).Error
    return member, err
}