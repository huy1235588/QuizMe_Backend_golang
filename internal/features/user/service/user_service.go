package service

import (
	"errors"

	"github.com/huy/quizme-backend/internal/features/user/dto"
	"github.com/huy/quizme-backend/internal/features/user/repository"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

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

func (s *userService) GetUser(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return dto.FromUser(user), nil
}

func (s *userService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		resp := dto.FromUser(&user)
		if resp != nil {
			responses[i] = *resp
		}
	}
	return responses, nil
}

func (s *userService) GetUsersPaged(page, pageSize int, search, sortBy, sortDir string) ([]dto.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAllPaged(page, pageSize, search, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		resp := dto.FromUser(&user)
		if resp != nil {
			responses[i] = *resp
		}
	}
	return responses, total, nil
}

func (s *userService) GetTopUsers(limit int) ([]dto.UserResponse, error) {
	users, err := s.userRepo.FindTopByTotalQuizPlays(limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		resp := dto.FromUser(&user)
		if resp != nil {
			responses[i] = *resp
		}
	}
	return responses, nil
}

func (s *userService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Update fields
	if req.Username != nil && *req.Username != "" {
		user.Username = *req.Username
	}
	if req.Email != nil && *req.Email != "" {
		user.Email = *req.Email
	}
	if req.FullName != nil && *req.FullName != "" {
		user.FullName = *req.FullName
	}
	if req.ProfileImage != nil {
		user.ProfileImage = req.ProfileImage
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return dto.FromUser(user), nil
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

// GetUserByID retrieves a user by ID (alias for GetUser)
func (s *userService) GetUserByID(id uint) (*dto.UserResponse, error) {
	return s.GetUser(id)
}

// GetUserCount returns the total count of users
func (s *userService) GetUserCount() (int64, error) {
	return s.userRepo.Count()
}

// GetUserProfile retrieves a user's profile by ID (alias for GetUser)
func (s *userService) GetUserProfile(id uint) (*dto.UserResponse, error) {
	return s.GetUser(id)
}

// UpdateUserAvatar updates a user's avatar URL
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

// RemoveUserAvatar removes a user's avatar
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

// ToggleUserActiveStatus toggles a user's active status
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
