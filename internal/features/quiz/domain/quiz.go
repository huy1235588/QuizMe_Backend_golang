package domain

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain/enums"
	category "github.com/huy/quizme-backend/internal/features/category/domain"
	user "github.com/huy/quizme-backend/internal/features/user/domain"
)

// Quiz represents a quiz
type Quiz struct {
	ID             uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	Title          string            `gorm:"size:100;not null" json:"title"`
	Description    *string           `gorm:"size:1000" json:"description,omitempty"`
	QuizThumbnails *string           `gorm:"column:quiz_thumbnails;size:255" json:"quizThumbnails,omitempty"`
	CreatorID      uint              `gorm:"column:creator_id;not null;index" json:"creatorId"`
	Difficulty     enums.Difficulty  `gorm:"type:varchar(20);not null;default:'MEDIUM';index" json:"difficulty"`
	IsPublic       bool              `gorm:"column:is_public;not null;default:true;index" json:"isPublic"`
	PlayCount      int               `gorm:"column:play_count;not null;default:0;index" json:"playCount"`
	QuestionCount  int               `gorm:"column:question_count;not null;default:0" json:"questionCount"`
	FavoriteCount  int               `gorm:"column:favorite_count;not null;default:0" json:"favoriteCount"`
	CreatedAt      time.Time         `gorm:"column:created_at;autoCreateTime;not null" json:"createdAt"`
	UpdatedAt      time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// Associations
	Creator    *user.User          `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Categories []category.Category `gorm:"many2many:quiz_category" json:"categories,omitempty"`
	Questions  []Question          `gorm:"foreignKey:QuizID" json:"questions,omitempty"`
}

// TableName specifies the table name for Quiz
func (Quiz) TableName() string {
	return "quiz"
}
