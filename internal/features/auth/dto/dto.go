package dto

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/features/user/domain"
)

// ===== REQUEST DTOs =====

// LoginRequest represents the login request body
type LoginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=2,max=50"`
	Email           string `json:"email" validate:"required,email,max=100"`
	Password        string `json:"password" validate:"required,min=2,max=100"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=2,max=100"`
	FullName        string `json:"fullName" validate:"required,max=100"`
}

// TokenRequest represents the refresh token request body
type TokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// ===== RESPONSE DTOs =====

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken        string        `json:"accessToken"`
	AccessTokenExpiry  time.Time     `json:"accessTokenExpiry"`
	RefreshToken       string        `json:"refreshToken"`
	RefreshTokenExpiry time.Time     `json:"refreshTokenExpiry"`
	User               *UserResponse `json:"user"`
}

// NewAuthResponse creates a new AuthResponse
func NewAuthResponse(
	accessToken string,
	accessTokenExpiry time.Time,
	refreshToken string,
	refreshTokenExpiry time.Time,
	user *UserResponse,
) *AuthResponse {
	return &AuthResponse{
		AccessToken:        accessToken,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshTokenExpiry,
		User:               user,
	}
}

// AccessTokenResponse represents a response with only access token (for refresh)
type AccessTokenResponse struct {
	AccessToken       string    `json:"accessToken"`
	AccessTokenExpiry time.Time `json:"accessTokenExpiry"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	FullName     string     `json:"fullName"`
	ProfileImage *string    `json:"profileImage,omitempty"`
	CreatedAt    string     `json:"createdAt"`
	UpdatedAt    string     `json:"updatedAt"`
	LastLogin    *string    `json:"lastLogin,omitempty"`
	Role         enums.Role `json:"role"`
	IsActive     bool       `json:"isActive"`
}

// FromUser converts a User domain model to UserResponse
func FromUser(user *domain.User) *UserResponse {
	var lastLogin *string
	if user.LastLogin != nil {
		l := user.LastLogin.Format(time.RFC3339)
		lastLogin = &l
	}

	return &UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		FullName:     user.FullName,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		LastLogin:    lastLogin,
		Role:         user.Role,
		IsActive:     user.IsActive,
	}
}

// FromUserWithImageURL converts a User domain model to UserResponse with image URL
func FromUserWithImageURL(user *domain.User, imageURL string) *UserResponse {
	response := FromUser(user)
	if imageURL != "" {
		response.ProfileImage = &imageURL
	}
	return response
}
