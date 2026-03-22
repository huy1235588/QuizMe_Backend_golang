package repository

import (
	"github.com/huy/quizme-backend/internal/domain"
	"gorm.io/gorm"
)

// GamePlayerAnswerRepository handles game player answer data access
type GamePlayerAnswerRepository interface {
	Create(answer *domain.GamePlayerAnswer) error
	CreateBatch(answers []*domain.GamePlayerAnswer) error
	FindByGameResultID(gameResultID uint) ([]*domain.GamePlayerAnswer, error)
	FindByParticipantAndQuestion(participantID, questionID uint) (*domain.GamePlayerAnswer, error)
}

type gamePlayerAnswerRepository struct {
	db *gorm.DB
}

// NewGamePlayerAnswerRepository creates a new game player answer repository
func NewGamePlayerAnswerRepository(db *gorm.DB) GamePlayerAnswerRepository {
	return &gamePlayerAnswerRepository{db: db}
}

func (r *gamePlayerAnswerRepository) Create(answer *domain.GamePlayerAnswer) error {
	return r.db.Create(answer).Error
}

func (r *gamePlayerAnswerRepository) CreateBatch(answers []*domain.GamePlayerAnswer) error {
	if len(answers) == 0 {
		return nil
	}
	return r.db.Create(&answers).Error
}

func (r *gamePlayerAnswerRepository) FindByGameResultID(gameResultID uint) ([]*domain.GamePlayerAnswer, error) {
	var answers []*domain.GamePlayerAnswer
	err := r.db.Where("game_result_id = ?", gameResultID).
		Preload("SelectedOptions").
		Find(&answers).Error
	return answers, err
}

func (r *gamePlayerAnswerRepository) FindByParticipantAndQuestion(participantID, questionID uint) (*domain.GamePlayerAnswer, error) {
	var answer domain.GamePlayerAnswer
	err := r.db.Where("participant_id = ? AND question_id = ?", participantID, questionID).
		Preload("SelectedOptions").
		First(&answer).Error
	return &answer, err
}
