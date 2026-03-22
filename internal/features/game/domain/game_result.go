package domain

import (
	"time"

	"gorm.io/gorm"
	quiz "github.com/huy/quizme-backend/internal/features/quiz/domain"
	room "github.com/huy/quizme-backend/internal/features/room/domain"
)

// GameResult stores the overall result of a game session
type GameResult struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	RoomID           uint           `gorm:"not null;index" json:"roomId"`
	Room             *room.Room     `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	QuizID           uint           `gorm:"not null;index" json:"quizId"`
	Quiz             *quiz.Quiz     `gorm:"foreignKey:QuizID" json:"quiz,omitempty"`
	StartTime        time.Time      `gorm:"index" json:"startTime"`
	EndTime          *time.Time     `json:"endTime,omitempty"`
	ParticipantCount int            `gorm:"not null" json:"participantCount"`
	QuestionCount    int            `gorm:"not null" json:"questionCount"`
	AvgScore         *float64       `json:"avgScore,omitempty"`
	HighestScore     *int           `json:"highestScore,omitempty"`
	LowestScore      *int           `json:"lowestScore,omitempty"`
	CompletionRate   *float64       `json:"completionRate,omitempty"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	GameResultQuestions []GameResultQuestion `gorm:"foreignKey:GameResultID" json:"gameResultQuestions,omitempty"`
	GamePlayerAnswers   []GamePlayerAnswer   `gorm:"foreignKey:GameResultID" json:"gamePlayerAnswers,omitempty"`
}

func (GameResult) TableName() string {
	return "game_results"
}

// GameResultQuestion stores statistics for each question in a game
type GameResultQuestion struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	GameResultID   uint        `gorm:"not null;index" json:"gameResultId"`
	GameResult     *GameResult `gorm:"foreignKey:GameResultID" json:"-"`
	QuestionID     uint        `gorm:"not null" json:"questionId"`
	Question       *quiz.Question `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	CorrectCount   int         `gorm:"not null;default:0" json:"correctCount"`
	IncorrectCount int         `gorm:"not null;default:0" json:"incorrectCount"`
	AvgTime        *float64    `json:"avgTime,omitempty"`
}

func (GameResultQuestion) TableName() string {
	return "game_result_questions"
}
