package repository

import "github.com/huy/quizme-backend/internal/features/game/domain"

// GameResultRepository defines the interface for game result data access
type GameResultRepository interface {
	Create(result *domain.GameResult) error
	FindByID(id uint) (*domain.GameResult, error)
	FindByRoomID(roomID uint) ([]*domain.GameResult, error)
	Update(result *domain.GameResult) error
}

// GamePlayerAnswerRepository defines the interface for game player answer data access
type GamePlayerAnswerRepository interface {
	Create(answer *domain.GamePlayerAnswer) error
	CreateBatch(answers []*domain.GamePlayerAnswer) error
	FindByGameResultID(gameResultID uint) ([]*domain.GamePlayerAnswer, error)
	FindByParticipantAndQuestion(participantID, questionID uint) (*domain.GamePlayerAnswer, error)
}
