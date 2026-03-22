package domain

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain/enums"
)

// Question represents a quiz question
type Question struct {
	ID          uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	QuizID      uint               `gorm:"column:quiz_id;not null;index" json:"quizId"`
	Content     string             `gorm:"type:text;not null" json:"content"`
	ImageURL    *string            `gorm:"column:image_url;size:255" json:"imageUrl,omitempty"`
	VideoURL    *string            `gorm:"column:video_url;size:255" json:"videoUrl,omitempty"`
	AudioURL    *string            `gorm:"column:audio_url;size:255" json:"audioUrl,omitempty"`
	FunFact     *string            `gorm:"column:fun_fact;type:text" json:"funFact,omitempty"`
	Explanation *string            `gorm:"type:text" json:"explanation,omitempty"`
	TimeLimit   int                `gorm:"column:time_limit;not null;default:30" json:"timeLimit"`
	Points      int                `gorm:"not null;default:10" json:"points"`
	OrderNumber int                `gorm:"column:order_number;not null;index" json:"orderNumber"`
	Type        enums.QuestionType `gorm:"column:type;type:varchar(20);not null;default:'QUIZ'" json:"type"`
	CreatedAt   time.Time          `gorm:"column:created_at;autoCreateTime;not null" json:"createdAt"`
	UpdatedAt   time.Time          `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// Associations
	Quiz    *Quiz            `gorm:"foreignKey:QuizID" json:"-"`
	Options []QuestionOption `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
}

// TableName specifies the table name for Question
func (Question) TableName() string {
	return "question"
}

// QuestionOption represents an answer option for a question
type QuestionOption struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID uint      `gorm:"column:question_id;not null;index" json:"questionId"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsCorrect  bool      `gorm:"column:is_correct;not null;default:false" json:"isCorrect"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;not null" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// Association
	Question *Question `gorm:"foreignKey:QuestionID" json:"-"`
}

// TableName specifies the table name for QuestionOption
func (QuestionOption) TableName() string {
	return "question_option"
}
