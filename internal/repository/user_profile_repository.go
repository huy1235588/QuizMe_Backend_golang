package repository

import (
	"github.com/huy/quizme-backend/internal/domain"
	"gorm.io/gorm"
)

// UserProfileRepository defines the interface for user profile data access
type UserProfileRepository interface {
	Create(profile *domain.UserProfile) error
	FindByID(id uint) (*domain.UserProfile, error)
	FindByUserID(userID uint) (*domain.UserProfile, error)
	Update(profile *domain.UserProfile) error
	Delete(id uint) error
}

type userProfileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository creates a new user profile repository
func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(profile *domain.UserProfile) error {
	return r.db.Create(profile).Error
}

func (r *userProfileRepository) FindByID(id uint) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	err := r.db.First(&profile, id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userProfileRepository) FindByUserID(userID uint) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userProfileRepository) Update(profile *domain.UserProfile) error {
	return r.db.Save(profile).Error
}

func (r *userProfileRepository) Delete(id uint) error {
	return r.db.Delete(&domain.UserProfile{}, id).Error
}
