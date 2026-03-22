package domain

import (
	"time"
)

// Category represents a quiz category
type Category struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description    *string   `gorm:"type:text" json:"description,omitempty"`
	IconURL        *string   `gorm:"column:icon_url;size:255" json:"iconUrl,omitempty"`
	QuizCount      int       `gorm:"column:quiz_count;not null;default:0" json:"quizCount"`
	TotalPlayCount int       `gorm:"column:total_play_count;not null;default:0" json:"totalPlayCount"`
	IsActive       bool      `gorm:"column:is_active;not null;default:true" json:"isActive"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime;not null" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// Associations (many-to-many with Quiz)
	Quizzes []Quiz `gorm:"many2many:quiz_category" json:"quizzes,omitempty"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "category"
}
