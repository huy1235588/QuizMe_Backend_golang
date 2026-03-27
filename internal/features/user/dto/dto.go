package dto

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/features/user/domain"
)

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

// UpdateUserRequest represents a request to update user information
type UpdateUserRequest struct {
	Username     *string     `json:"username,omitempty"`
	Email        *string     `json:"email,omitempty"`
	FullName     *string     `json:"fullName,omitempty"`
	Password     *string     `json:"password,omitempty"`
	ProfileImage *string     `json:"profileImage,omitempty"`
	Role         *enums.Role `json:"role,omitempty"`
	IsActive     *bool       `json:"isActive,omitempty"`
}

// CreateUserRequest represents a request to create a new user (admin only)
type CreateUserRequest struct {
	Username string      `json:"username" validate:"required,min=2,max=50"`
	Email    string      `json:"email" validate:"required,email,max=100"`
	FullName string      `json:"fullName" validate:"required,max=100"`
	Password *string     `json:"password,omitempty" validate:"omitempty,min=6,max=100"`
	Role     *enums.Role `json:"role,omitempty"`
	IsActive *bool       `json:"isActive,omitempty"`
}

// UserProfileResponse represents a user profile in API responses
type UserProfileResponse struct {
	ID             uint    `json:"id"`
	UserID         uint    `json:"userId"`
	DateOfBirth    *string `json:"dateOfBirth,omitempty"`
	City           *string `json:"city,omitempty"`
	PhoneNumber    *string `json:"phoneNumber,omitempty"`
	TotalScore     int     `json:"totalScore"`
	QuizzesPlayed  int     `json:"quizzesPlayed"`
	QuizzesCreated int     `json:"quizzesCreated"`
	TotalQuizPlays int     `json:"totalQuizPlays"`
}

// FromUserProfile converts a UserProfile domain model to UserProfileResponse
func FromUserProfile(profile *domain.UserProfile) *UserProfileResponse {
	var dateOfBirth *string
	if profile.DateOfBirth != nil {
		dob := profile.DateOfBirth.Format("2006-01-02")
		dateOfBirth = &dob
	}

	return &UserProfileResponse{
		ID:             profile.ID,
		UserID:         profile.UserID,
		DateOfBirth:    dateOfBirth,
		City:           profile.City,
		PhoneNumber:    profile.PhoneNumber,
		TotalScore:     profile.TotalScore,
		QuizzesPlayed:  profile.QuizzesPlayed,
		QuizzesCreated: profile.QuizzesCreated,
		TotalQuizPlays: profile.TotalQuizPlays,
	}
}
