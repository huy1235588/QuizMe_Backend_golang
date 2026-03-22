package repository

import (
	"github.com/huy/quizme-backend/internal/features/game/domain"
	"gorm.io/gorm"
)

type gameResultRepository struct {
	db *gorm.DB
}

// NewGameResultRepository creates a new game result repository
func NewGameResultRepository(db *gorm.DB) GameResultRepository {
	return &gameResultRepository{db: db}
}

func (r *gameResultRepository) Create(result *domain.GameResult) error {
	return r.db.Create(result).Error
}

func (r *gameResultRepository) FindByID(id uint) (*domain.GameResult, error) {
	var result domain.GameResult
	err := r.db.Preload("GameResultQuestions").Preload("GamePlayerAnswers").First(&result, id).Error
	return &result, err
}

func (r *gameResultRepository) FindByRoomID(roomID uint) ([]*domain.GameResult, error) {
	var results []*domain.GameResult
	err := r.db.Where("room_id = ?", roomID).Order("created_at desc").Find(&results).Error
	return results, err
}

func (r *gameResultRepository) Update(result *domain.GameResult) error {
	return r.db.Save(result).Error
}
