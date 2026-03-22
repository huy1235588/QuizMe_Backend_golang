package service

import (
	"errors"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserService handles user-related operations
type UserService interface {
	GetUserByID(id uint) (*response.UserResponse, error)
	GetUserProfile(id uint) (*domain.User, error)
	GetTopUsers(limit int) ([]*response.UserResponse, error)
	GetUserCount() (int64, error)
	GetPagedUsers(page, pageSize int, search, sortBy, sortDir string) ([]*response.UserResponse, int64, error)
	UpdateUser(id uint, user *domain.User) error
	DeleteUser(id uint) error
	ToggleUserActiveStatus(id uint, isActive bool) error
	UpdateUserAvatar(id uint, avatarURL string) error
	RemoveUserAvatar(id uint) error
}

type userService struct {
	userRepo        repository.UserRepository
	userProfileRepo repository.UserProfileRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, userProfileRepo repository.UserProfileRepository) UserService {
	return &userService{
		userRepo:        userRepo,
		userProfileRepo: userProfileRepo,
	}
}

func (s *userService) GetUserByID(id uint) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return response.FromUser(user), nil
}

func (s *userService) GetUserProfile(id uint) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetTopUsers(limit int) ([]*response.UserResponse, error) {
	users, err := s.userRepo.FindTopByTotalQuizPlays(limit)
	if err != nil {
		return nil, err
	}

	responses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		responses[i] = response.FromUser(&user)
	}
	return responses, nil
}

func (s *userService) GetUserCount() (int64, error) {
	return s.userRepo.Count()
}

func (s *userService) GetPagedUsers(page, pageSize int, search, sortBy, sortDir string) ([]*response.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAllPaged(page, pageSize, search, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		responses[i] = response.FromUser(&user)
	}
	return responses, total, nil
}

func (s *userService) UpdateUser(id uint, updatedUser *domain.User) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Update fields
	if updatedUser.Username != "" {
		user.Username = updatedUser.Username
	}
	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}
	if updatedUser.FullName != "" {
		user.FullName = updatedUser.FullName
	}
	if updatedUser.Password != "" {
		user.Password = updatedUser.Password
	}
	if updatedUser.ProfileImage != nil {
		user.ProfileImage = updatedUser.ProfileImage
	}
	user.Role = updatedUser.Role
	user.IsActive = updatedUser.IsActive

	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return s.userRepo.Delete(id)
}

func (s *userService) ToggleUserActiveStatus(id uint, isActive bool) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	user.IsActive = isActive
	return s.userRepo.Update(user)
}

func (s *userService) UpdateUserAvatar(id uint, avatarURL string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	user.ProfileImage = &avatarURL
	return s.userRepo.Update(user)
}

func (s *userService) RemoveUserAvatar(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	user.ProfileImage = nil
	return s.userRepo.Update(user)
}
