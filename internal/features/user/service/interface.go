package service

import (
	"github.com/huy/quizme-backend/internal/features/user/dto"
)

// UserService defines the interface for user operations
type UserService interface {
	GetUser(id uint) (*dto.UserResponse, error)
	GetAllUsers() ([]dto.UserResponse, error)
	GetUsersPaged(page, pageSize int, search, sortBy, sortDir string) ([]dto.UserResponse, int64, error)
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uint) error
	GetTopUsers(limit int) ([]dto.UserResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	GetUserCount() (int64, error)
	GetUserProfile(id uint) (*dto.UserProfileResponse, error)
	UpdateUserAvatar(id uint, avatarURL string) error
	RemoveUserAvatar(id uint) error
	ToggleUserActiveStatus(id uint, isActive bool) error
}
