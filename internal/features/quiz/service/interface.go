package service

import (
	"github.com/huy/quizme-backend/internal/domain/enums"
	quizdto "github.com/huy/quizme-backend/internal/features/quiz/dto"
)

// QuizService handles quiz-related operations
type QuizService interface {
	GetAllQuizzes() ([]*quizdto.QuizResponse, error)
	GetPublicQuizzes() ([]*quizdto.QuizResponse, error)
	GetQuizByID(id uint) (*quizdto.QuizResponse, error)
	GetQuizWithQuestions(id uint) (*quizdto.QuizResponse, []quizdto.QuestionResponse, error)
	GetQuizzesByDifficulty(difficulty enums.Difficulty) ([]*quizdto.QuizResponse, error)
	GetQuizzesWithFilters(categoryID *uint, difficulty *string, isPublic *bool, search *string, page, pageSize int, sortBy, sortDir string) ([]*quizdto.QuizResponse, int64, error)
	CreateQuiz(creatorID uint, req *quizdto.QuizRequest) (*quizdto.QuizResponse, error)
	UpdateQuiz(id, userID uint, req *quizdto.QuizRequest) (*quizdto.QuizResponse, error)
	DeleteQuiz(id, userID uint) error
}

// QuestionService handles question-related operations
type QuestionService interface {
	GetAllQuestions() ([]*quizdto.QuestionResponse, error)
	GetQuestionByID(id uint) (*quizdto.QuestionResponse, error)
	GetQuestionsByQuizID(quizID uint) ([]*quizdto.QuestionResponse, error)
	CreateQuestion(req *quizdto.QuestionRequest, quizID uint) (*quizdto.QuestionResponse, error)
	CreateBatchQuestions(req *quizdto.BatchQuestionRequest) ([]*quizdto.QuestionResponse, error)
	UpdateQuestion(id uint, req *quizdto.QuestionRequest) (*quizdto.QuestionResponse, error)
	DeleteQuestion(id uint) error
}
