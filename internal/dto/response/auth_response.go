package response

import (
	"time"
)

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
