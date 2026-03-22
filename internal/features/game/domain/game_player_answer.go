package domain

import (
	"time"

	"gorm.io/gorm"
	quiz "github.com/huy/quizme-backend/internal/features/quiz/domain"
	room "github.com/huy/quizme-backend/internal/features/room/domain"
)

// GamePlayerAnswer stores individual player answers during a game
type GamePlayerAnswer struct {
	ID            uint                  `gorm:"primaryKey" json:"id"`
	GameResultID  uint                  `gorm:"not null;index:idx_game_result_participant" json:"gameResultId"`
	GameResult    *GameResult           `gorm:"foreignKey:GameResultID" json:"-"`
	ParticipantID uint                  `gorm:"not null;index:idx_game_result_participant;index:idx_participant_question" json:"participantId"`
	Participant   *room.RoomParticipant `gorm:"foreignKey:ParticipantID" json:"participant,omitempty"`
	QuestionID    uint                  `gorm:"not null;index:idx_participant_question" json:"questionId"`
	Question      *quiz.Question        `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	IsCorrect     bool                  `gorm:"default:false" json:"isCorrect"`
	AnswerTime    float64               `gorm:"not null" json:"answerTime"`
	Score         int                   `gorm:"not null;default:0" json:"score"`
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt        `gorm:"index" json:"-"`

	// Relations
	SelectedOptions []GamePlayerAnswerOption `gorm:"foreignKey:GamePlayerAnswerID" json:"selectedOptions,omitempty"`
}

func (GamePlayerAnswer) TableName() string {
	return "game_player_answers"
}

// GetSelectedOptionIDs returns the IDs of selected options
func (a *GamePlayerAnswer) GetSelectedOptionIDs() []uint {
	ids := make([]uint, len(a.SelectedOptions))
	for i, opt := range a.SelectedOptions {
		ids[i] = opt.OptionID
	}
	return ids
}

// GamePlayerAnswerOption stores the selected options for a player answer
type GamePlayerAnswerOption struct {
	ID                 uint              `gorm:"primaryKey" json:"id"`
	GamePlayerAnswerID uint              `gorm:"not null;index;uniqueIndex:unique_player_answer_option" json:"gamePlayerAnswerId"`
	GamePlayerAnswer   *GamePlayerAnswer `gorm:"foreignKey:GamePlayerAnswerID" json:"-"`
	OptionID           uint              `gorm:"not null;uniqueIndex:unique_player_answer_option" json:"optionId"`
	Option             *quiz.QuestionOption `gorm:"foreignKey:OptionID" json:"option,omitempty"`
}

func (GamePlayerAnswerOption) TableName() string {
	return "game_player_answer_options"
}
