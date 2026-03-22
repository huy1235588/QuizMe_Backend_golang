package service

import (
	"errors"
	"time"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/dto/request"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/pkg/jwt"
	"github.com/huy/quizme-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials   = errors.New("invalid username/email or password")
	ErrUsernameExists       = errors.New("username already exists")
	ErrEmailExists          = errors.New("email already exists")
	ErrPasswordMismatch     = errors.New("passwords do not match")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrExpiredRefreshToken  = errors.New("refresh token expired")
	ErrUserNotActive        = errors.New("user account is not active")
)

// AuthService handles authentication operations
type AuthService interface {
	Login(req *request.LoginRequest) (*response.AuthResponse, error)
	Register(req *request.RegisterRequest) (*response.AuthResponse, error)
	Logout(refreshToken string) error
	RefreshToken(refreshToken string) (*response.AuthResponse, error)
}

type authService struct {
	userRepo         repository.UserRepository
	userProfileRepo  repository.UserProfileRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtProvider      *jwt.JWTProvider
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repository.UserRepository,
	userProfileRepo repository.UserProfileRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtProvider *jwt.JWTProvider,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		userProfileRepo:  userProfileRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtProvider:      jwtProvider,
	}
}

// Login authenticates a user and returns tokens
func (s *authService) Login(req *request.LoginRequest) (*response.AuthResponse, error) {
	// Find user by username or email
	user, err := s.userRepo.FindByUsernameOrEmail(req.UsernameOrEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	return s.createAuthResponse(user)
}

// Register creates a new user and returns tokens
func (s *authService) Register(req *request.RegisterRequest) (*response.AuthResponse, error) {
	// Check if username exists
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	// Check if email exists
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Verify passwords match
	if req.Password != req.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Role:     enums.RoleUser,
		IsActive: true,
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

	// Generate tokens
	return s.createAuthResponse(user)
}

// Logout invalidates the refresh token
func (s *authService) Logout(refreshToken string) error {
	return s.refreshTokenRepo.DeleteByToken(refreshToken)
}

// RefreshToken generates a new access token using a valid refresh token
func (s *authService) RefreshToken(refreshToken string) (*response.AuthResponse, error) {
	// Find refresh token
	token, err := s.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	// Check if token is expired
	if token.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredRefreshToken
	}

	// Check if token is revoked
	if token.Revoked {
		return nil, ErrInvalidRefreshToken
	}

	// Find user
	user, err := s.userRepo.FindByID(token.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new access token only
	accessToken, accessExpiry, err := s.jwtProvider.GenerateAccessToken(user.Username)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken:        accessToken,
		AccessTokenExpiry:  accessExpiry,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: token.ExpiresAt,
		User:               response.FromUser(user),
	}, nil
}

// createAuthResponse generates tokens and creates the auth response
func (s *authService) createAuthResponse(user *domain.User) (*response.AuthResponse, error) {
	// Generate access token
	accessToken, accessExpiry, err := s.jwtProvider.GenerateAccessToken(user.Username)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, jti, refreshExpiry, err := s.jwtProvider.GenerateRefreshToken(user.Username)
	if err != nil {
		return nil, err
	}

	// Delete old refresh tokens for this user
	_ = s.refreshTokenRepo.DeleteByUserID(user.ID)

	// Save new refresh token
	tokenEntity := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		JTI:       jti,
		IssuedAt:  time.Now(),
		ExpiresAt: refreshExpiry,
		Revoked:   false,
	}

	if err := s.refreshTokenRepo.Create(tokenEntity); err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken:        accessToken,
		AccessTokenExpiry:  accessExpiry,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshExpiry,
		User:               response.FromUser(user),
	}, nil
}
