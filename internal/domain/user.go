package domain

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain/enums"
)

// User represents a user in the system
type User struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string     `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password     string     `gorm:"not null" json:"-"` // Never expose password in JSON
	FullName     string     `gorm:"column:full_name;size:100;not null" json:"fullName"`
	ProfileImage *string    `gorm:"column:profile_image;size:255" json:"profileImage,omitempty"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime;not null" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime;not null" json:"updatedAt"`
	LastLogin    *time.Time `gorm:"column:last_login" json:"lastLogin,omitempty"`
	Role         enums.Role `gorm:"column:role;type:varchar(20);not null;default:'USER'" json:"role"`
	IsActive     bool       `gorm:"column:is_active;not null;default:true" json:"isActive"`

	// Associations
	UserProfile *UserProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"userProfile,omitempty"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "user"
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == enums.RoleAdmin
}

// IsEnabled checks if the user account is enabled
func (u *User) IsEnabled() bool {
	return u.IsActive
}
