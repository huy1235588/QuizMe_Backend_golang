package service

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/features/user/domain"
	"github.com/huy/quizme-backend/internal/features/user/dto"
	"github.com/huy/quizme-backend/internal/features/user/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUsernameExists      = errors.New("username already exists")
	ErrEmailExists         = errors.New("email already exists")
	ErrInvalidRole         = errors.New("invalid role specified")
	ErrUserProfileNotFound = errors.New("user profile not found")
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

// CreateUser creates a new user (admin only)
func (s *userService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check username uniqueness
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	// Check email uniqueness
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Handle password
	var password string
	if req.Password != nil && *req.Password != "" {
		password = *req.Password
	} else {
		password = generateRandomPassword(12)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Determine role (default to USER)
	role := enums.RoleUser
	if req.Role != nil {
		if !req.Role.IsValid() {
			return nil, ErrInvalidRole
		}
		role = *req.Role
	}

	// Determine isActive (default to true)
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Create user
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Role:     role,
		IsActive: isActive,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Create user profile
	profile := &domain.UserProfile{
		UserID: user.ID,
	}
	if err := s.userProfileRepo.Create(profile); err != nil {
		return nil, err
	}

	return dto.FromUser(user), nil
}

// generateRandomPassword generates a secure random password
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func (s *userService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Update username with uniqueness check
	if req.Username != nil && *req.Username != "" && *req.Username != user.Username {
		exists, err := s.userRepo.ExistsByUsername(*req.Username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrUsernameExists
		}
		user.Username = *req.Username
	}

	// Update email with uniqueness check
	if req.Email != nil && *req.Email != "" && *req.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(*req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailExists
		}
		user.Email = *req.Email
	}

	// Update full name
	if req.FullName != nil && *req.FullName != "" {
		user.FullName = *req.FullName
	}

	// Update and hash password if provided
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	// Update role if provided
	if req.Role != nil {
		if !req.Role.IsValid() {
			return nil, ErrInvalidRole
		}
		user.Role = *req.Role
	}

	// Update profile image
	if req.ProfileImage != nil {
		user.ProfileImage = req.ProfileImage
	}

	// Update active status
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

// GetUserProfile retrieves a user's profile by user ID
func (s *userService) GetUserProfile(id uint) (*dto.UserProfileResponse, error) {
	profile, err := s.userProfileRepo.FindByUserID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserProfileNotFound
		}
		return nil, err
	}
	return dto.FromUserProfile(profile), nil
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
