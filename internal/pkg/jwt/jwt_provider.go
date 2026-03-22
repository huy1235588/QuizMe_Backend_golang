package jwt

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTProvider handles JWT token generation and validation
type JWTProvider struct {
	secretKey             []byte
	accessExpirationMs    int64
	refreshExpirationMs   int64
}

// Claims represents the JWT claims
type Claims struct {
	jwt.RegisteredClaims
}

// NewJWTProvider creates a new JWT provider
func NewJWTProvider(secret string, accessExpirationMs, refreshExpirationMs int64) *JWTProvider {
	// Decode base64 secret (matching Spring Boot's Base64.decode)
	secretKey, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		// If not valid base64, use raw secret
		secretKey = []byte(secret)
	}

	return &JWTProvider{
		secretKey:           secretKey,
		accessExpirationMs:  accessExpirationMs,
		refreshExpirationMs: refreshExpirationMs,
	}
}

// GenerateAccessToken generates an access token for the given username
func (p *JWTProvider) GenerateAccessToken(username string) (string, time.Time, error) {
	return p.generateToken(username, p.accessExpirationMs)
}

// GenerateRefreshToken generates a refresh token for the given username
func (p *JWTProvider) GenerateRefreshToken(username string) (string, string, time.Time, error) {
	token, expiresAt, err := p.generateToken(username, p.refreshExpirationMs)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Extract JTI from the token
	jti, err := p.GetJTIFromToken(token)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return token, jti, expiresAt, nil
}

// generateToken creates a JWT token with the given subject and expiration
func (p *JWTProvider) generateToken(subject string, expirationMs int64) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(expirationMs) * time.Millisecond)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString(p.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// ValidateToken validates a JWT token and returns the claims
func (p *JWTProvider) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method is HS512
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return p.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetUsernameFromToken extracts the username (subject) from a token
func (p *JWTProvider) GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := p.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}

// GetJTIFromToken extracts the JWT ID from a token
func (p *JWTProvider) GetJTIFromToken(tokenString string) (string, error) {
	claims, err := p.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.ID, nil
}

// GetExpirationFromToken extracts the expiration time from a token
func (p *JWTProvider) GetExpirationFromToken(tokenString string) (time.Time, error) {
	claims, err := p.ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return claims.ExpiresAt.Time, nil
}

// GetAccessExpirationMs returns access token expiration in milliseconds
func (p *JWTProvider) GetAccessExpirationMs() int64 {
	return p.accessExpirationMs
}

// GetRefreshExpirationMs returns refresh token expiration in milliseconds
func (p *JWTProvider) GetRefreshExpirationMs() int64 {
	return p.refreshExpirationMs
}
