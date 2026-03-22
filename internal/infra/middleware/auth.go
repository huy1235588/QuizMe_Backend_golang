package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/huy/quizme-backend/internal/features/user/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/pkg/jwt"
	userrepo "github.com/huy/quizme-backend/internal/features/user/repository"
)

const (
	// AuthorizationHeader is the header key for authorization
	AuthorizationHeader = "Authorization"
	// BearerPrefix is the prefix for bearer tokens
	BearerPrefix = "Bearer "
	// CurrentUserKey is the context key for the current user
	CurrentUserKey = "currentUser"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtProvider *jwt.JWTProvider
	userRepo    userrepo.UserRepository
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtProvider *jwt.JWTProvider, userRepo userrepo.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtProvider: jwtProvider,
		userRepo:    userRepo,
	}
}

// RequireAuth returns a middleware that requires authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Error("Authorization header is required"))
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, response.Error("Invalid authorization header format"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)

		// Validate token
		claims, err := m.jwtProvider.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("Invalid or expired token"))
			c.Abort()
			return
		}

		// Get user from database
		user, err := m.userRepo.FindByUsername(claims.Subject)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("User not found"))
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			c.JSON(http.StatusForbidden, response.Error("User account is not active"))
			c.Abort()
			return
		}

		// Set user in context
		c.Set(CurrentUserKey, user)
		c.Next()
	}
}

// OptionalAuth returns a middleware that optionally authenticates
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" || !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)

		claims, err := m.jwtProvider.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		user, err := m.userRepo.FindByUsername(claims.Subject)
		if err != nil {
			c.Next()
			return
		}

		if user.IsActive {
			c.Set(CurrentUserKey, user)
		}

		c.Next()
	}
}

// RequireAdmin returns a middleware that requires admin role
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(CurrentUserKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
			c.Abort()
			return
		}

		currentUser := user.(*domain.User)
		if currentUser.Role != enums.RoleAdmin {
			c.JSON(http.StatusForbidden, response.Error("Access denied: Admin role required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUser retrieves the current user from context
func GetCurrentUser(c *gin.Context) *domain.User {
	user, exists := c.Get(CurrentUserKey)
	if !exists {
		return nil
	}
	return user.(*domain.User)
}

// ValidateToken validates a JWT token and returns the user
func (m *AuthMiddleware) ValidateToken(token string) (*domain.User, error) {
	claims, err := m.jwtProvider.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	user, err := m.userRepo.FindByUsername(claims.Subject)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, jwt.ErrInvalidToken
	}

	return user, nil
}
