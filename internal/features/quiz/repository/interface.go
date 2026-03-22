package repository

import (
	"github.com/huy/quizme-backend/internal/features/quiz/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
)

// QuizRepository defines the interface for quiz data access
type QuizRepository interface {
	Create(quiz *domain.Quiz) error
	FindByID(id uint) (*domain.Quiz, error)
	FindByIDWithQuestions(id uint) (*domain.Quiz, error)
	FindAll() ([]domain.Quiz, error)
	FindPublic() ([]domain.Quiz, error)
	FindByDifficulty(difficulty enums.Difficulty) ([]domain.Quiz, error)
	FindByCreatorID(creatorID uint) ([]domain.Quiz, error)
	FindWithFilters(categoryID *uint, difficulty *string, isPublic *bool, search *string, page, pageSize int, sortBy, sortDir string) ([]domain.Quiz, int64, error)
	Update(quiz *domain.Quiz) error
	Delete(id uint) error
	IncrementPlayCount(id uint) error
	UpdateQuestionCount(id uint, count int) error
	AddCategories(quizID uint, categoryIDs []uint) error
	RemoveCategories(quizID uint) error
}

// QuestionRepository defines the interface for question data access
type QuestionRepository interface {
	Create(question *domain.Question) error
	CreateBatch(questions []domain.Question) error
	FindByID(id uint) (*domain.Question, error)
	FindByQuizID(quizID uint) ([]domain.Question, error)
	Update(question *domain.Question) error
	Delete(id uint) error
	DeleteByQuizID(quizID uint) error
	CountByQuizID(quizID uint) (int64, error)
}

// QuestionOptionRepository defines the interface for question option data access
type QuestionOptionRepository interface {
	Create(option *domain.QuestionOption) error
	CreateBatch(options []domain.QuestionOption) error
	FindByID(id uint) (*domain.QuestionOption, error)
	FindByQuestionID(questionID uint) ([]domain.QuestionOption, error)
	Update(option *domain.QuestionOption) error
	Delete(id uint) error
	DeleteByQuestionID(questionID uint) error
}
