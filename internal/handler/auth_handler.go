package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huy/quizme-backend/internal/dto/request"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	authResponse, err := h.authService.Login(&req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, response.Error("Invalid username/email or password"))
		case errors.Is(err, service.ErrUserNotActive):
			c.JSON(http.StatusForbidden, response.Error("User account is not active"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("An error occurred during login"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(authResponse, "Login successful"))
}

// Register handles user registration
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	authResponse, err := h.authService.Register(&req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUsernameExists):
			c.JSON(http.StatusConflict, response.Error("Username already exists"))
		case errors.Is(err, service.ErrEmailExists):
			c.JSON(http.StatusConflict, response.Error("Email already exists"))
		case errors.Is(err, service.ErrPasswordMismatch):
			c.JSON(http.StatusBadRequest, response.Error("Passwords do not match"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("An error occurred during registration"))
		}
		return
	}

	c.JSON(http.StatusCreated, response.Success(authResponse, "Registration successful"))
}

// Logout handles user logout
// POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req request.TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid refresh token"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "Logout successful"))
}

// RefreshToken handles token refresh
// POST /api/auth/refresh-token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	authResponse, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefreshToken):
			c.JSON(http.StatusUnauthorized, response.Error("Invalid refresh token"))
		case errors.Is(err, service.ErrExpiredRefreshToken):
			c.JSON(http.StatusUnauthorized, response.Error("Refresh token expired"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("An error occurred during token refresh"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(authResponse, "Token refreshed successfully"))
}
