package repository

import (
	authdomain "github.com/huy/quizme-backend/internal/features/auth/domain"
	userdomain "github.com/huy/quizme-backend/internal/features/user/domain"
)

// UserRepository defines the interface for user data access in auth feature
type UserRepository interface {
	Create(user *userdomain.User) error
	FindByID(id uint) (*userdomain.User, error)
	FindByUsername(username string) (*userdomain.User, error)
	FindByEmail(email string) (*userdomain.User, error)
	FindByUsernameOrEmail(usernameOrEmail string) (*userdomain.User, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Update(user *userdomain.User) error
	Delete(id uint) error
	FindAll() ([]userdomain.User, error)
	FindAllPaged(page, pageSize int, search, sortBy, sortDir string) ([]userdomain.User, int64, error)
	Count() (int64, error)
	FindTopByTotalQuizPlays(limit int) ([]userdomain.User, error)
}

// UserProfileRepository defines the interface for user profile data access in auth feature
type UserProfileRepository interface {
	Create(profile *authdomain.UserProfile) error
	FindByUserID(userID uint) (*authdomain.UserProfile, error)
	Update(profile *authdomain.UserProfile) error
	Delete(userID uint) error
}

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	Create(token *authdomain.RefreshToken) error
	FindByToken(token string) (*authdomain.RefreshToken, error)
	FindByJTI(jti string) (*authdomain.RefreshToken, error)
	FindByUserID(userID uint) (*authdomain.RefreshToken, error)
	DeleteByUserID(userID uint) error
	DeleteByToken(token string) error
	RevokeByUserID(userID uint) error
}
