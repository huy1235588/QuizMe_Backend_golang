package service

import (
	"github.com/huy/quizme-backend/internal/features/auth/dto"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Logout(refreshToken string) error
	RefreshToken(refreshToken string) (*dto.AuthResponse, error)
}
