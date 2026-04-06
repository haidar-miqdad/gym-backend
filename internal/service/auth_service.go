package service

import (
	"errors"
	"gym-backend/internal/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(username, password string) (string, error)
}

type authService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db}
}

func (s *authService) Login(username, password string) (string, error) {
	var user domain.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("username atau password salah")
	}

	// Cek Password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("username atau password salah")
	}

	// Buat JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Berlaku 24 jam
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}