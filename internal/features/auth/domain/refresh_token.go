package domain

import (
	"time"

	"github.com/huy/quizme-backend/internal/features/user/domain"
)

// RefreshToken stores refresh tokens for JWT authentication
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null;index" json:"userId"`
	Token     string    `gorm:"uniqueIndex;size:500;not null" json:"-"`
	JTI       string    `gorm:"column:jti;uniqueIndex;size:100;not null" json:"-"` // JWT ID
	IssuedAt  time.Time `gorm:"column:issued_at;not null" json:"issuedAt"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expiresAt"`
	Revoked   bool      `gorm:"not null;default:false" json:"revoked"`

	// Association
	User *domain.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName specifies the table name for RefreshToken
func (RefreshToken) TableName() string {
	return "refresh_token"
}

// IsExpired checks if the refresh token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if the refresh token is valid (not revoked and not expired)
func (rt *RefreshToken) IsValid() bool {
	return !rt.Revoked && !rt.IsExpired()
}
