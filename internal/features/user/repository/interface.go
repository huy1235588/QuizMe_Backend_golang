package repository

import (
	"github.com/huy/quizme-backend/internal/features/user/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByUsernameOrEmail(usernameOrEmail string) (*domain.User, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Update(user *domain.User) error
	Delete(id uint) error
	FindAll() ([]domain.User, error)
	FindAllPaged(page, pageSize int, search, sortBy, sortDir string) ([]domain.User, int64, error)
	Count() (int64, error)
	FindTopByTotalQuizPlays(limit int) ([]domain.User, error)
}

// UserProfileRepository defines the interface for user profile data access
type UserProfileRepository interface {
	Create(profile *domain.UserProfile) error
	FindByID(id uint) (*domain.UserProfile, error)
	FindByUserID(userID uint) (*domain.UserProfile, error)
	Update(profile *domain.UserProfile) error
	Delete(id uint) error
}
