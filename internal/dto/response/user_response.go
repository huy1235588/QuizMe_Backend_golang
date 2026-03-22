package response

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
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
