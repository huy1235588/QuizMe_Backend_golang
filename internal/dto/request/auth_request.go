package request

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
