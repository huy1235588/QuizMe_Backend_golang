package service

import (
	"errors"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/dto/request"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrQuizNotFound     = errors.New("quiz not found")
	ErrNotQuizOwner     = errors.New("you are not the owner of this quiz")
	ErrQuestionNotFound = errors.New("question not found")
)

// QuizService handles quiz-related operations
type QuizService interface {
	GetAllQuizzes() ([]*response.QuizResponse, error)
	GetPublicQuizzes() ([]*response.QuizResponse, error)
	GetQuizByID(id uint) (*response.QuizResponse, error)
	GetQuizWithQuestions(id uint) (*response.QuizResponse, []response.QuestionResponse, error)
	GetQuizzesByDifficulty(difficulty enums.Difficulty) ([]*response.QuizResponse, error)
	GetQuizzesWithFilters(categoryID *uint, difficulty *string, isPublic *bool, search *string, page, pageSize int, sortBy, sortDir string) ([]*response.QuizResponse, int64, error)
	CreateQuiz(creatorID uint, req *request.QuizRequest) (*response.QuizResponse, error)
	UpdateQuiz(id, userID uint, req *request.QuizRequest) (*response.QuizResponse, error)
	DeleteQuiz(id, userID uint) error
}

type quizService struct {
	quizRepo           repository.QuizRepository
	questionRepo       repository.QuestionRepository
	questionOptionRepo repository.QuestionOptionRepository
	categoryRepo       repository.CategoryRepository
}

// NewQuizService creates a new quiz service
func NewQuizService(
	quizRepo repository.QuizRepository,
	questionRepo repository.QuestionRepository,
	questionOptionRepo repository.QuestionOptionRepository,
	categoryRepo repository.CategoryRepository,
) QuizService {
	return &quizService{
		quizRepo:           quizRepo,
		questionRepo:       questionRepo,
		questionOptionRepo: questionOptionRepo,
		categoryRepo:       categoryRepo,
	}
}

func (s *quizService) GetAllQuizzes() ([]*response.QuizResponse, error) {
	quizzes, err := s.quizRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return response.FromQuizzes(quizzes), nil
}

func (s *quizService) GetPublicQuizzes() ([]*response.QuizResponse, error) {
	quizzes, err := s.quizRepo.FindPublic()
	if err != nil {
		return nil, err
	}
	return response.FromQuizzes(quizzes), nil
}

func (s *quizService) GetQuizByID(id uint) (*response.QuizResponse, error) {
	quiz, err := s.quizRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuizNotFound
		}
		return nil, err
	}
	return response.FromQuiz(quiz), nil
}

func (s *quizService) GetQuizWithQuestions(id uint) (*response.QuizResponse, []response.QuestionResponse, error) {
	quiz, err := s.quizRepo.FindByIDWithQuestions(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrQuizNotFound
		}
		return nil, nil, err
	}

	quizResponse := response.FromQuiz(quiz)
	questionResponses := make([]response.QuestionResponse, len(quiz.Questions))
	for i, q := range quiz.Questions {
		questionResponses[i] = *response.FromQuestion(&q)
	}

	return quizResponse, questionResponses, nil
}

func (s *quizService) GetQuizzesByDifficulty(difficulty enums.Difficulty) ([]*response.QuizResponse, error) {
	quizzes, err := s.quizRepo.FindByDifficulty(difficulty)
	if err != nil {
		return nil, err
	}
	return response.FromQuizzes(quizzes), nil
}

func (s *quizService) GetQuizzesWithFilters(
	categoryID *uint,
	difficulty *string,
	isPublic *bool,
	search *string,
	page, pageSize int,
	sortBy, sortDir string,
) ([]*response.QuizResponse, int64, error) {
	quizzes, total, err := s.quizRepo.FindWithFilters(categoryID, difficulty, isPublic, search, page, pageSize, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}
	return response.FromQuizzes(quizzes), total, nil
}

func (s *quizService) CreateQuiz(creatorID uint, req *request.QuizRequest) (*response.QuizResponse, error) {
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	quiz := &domain.Quiz{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   creatorID,
		Difficulty:  req.Difficulty,
		IsPublic:    isPublic,
	}

	if err := s.quizRepo.Create(quiz); err != nil {
		return nil, err
	}

	// Add categories
	if len(req.CategoryIDs) > 0 {
		if err := s.quizRepo.AddCategories(quiz.ID, req.CategoryIDs); err != nil {
			return nil, err
		}
		// Increment quiz count for categories
		for _, catID := range req.CategoryIDs {
			_ = s.categoryRepo.IncrementQuizCount(catID)
		}
	}

	// Create questions if provided
	if len(req.Questions) > 0 {
		for i, qReq := range req.Questions {
			question := &domain.Question{
				QuizID:      quiz.ID,
				Content:     qReq.Content,
				ImageURL:    qReq.ImageURL,
				VideoURL:    qReq.VideoURL,
				AudioURL:    qReq.AudioURL,
				FunFact:     qReq.FunFact,
				Explanation: qReq.Explanation,
				TimeLimit:   qReq.TimeLimit,
				Points:      qReq.Points,
				OrderNumber: i + 1,
				Type:        qReq.Type,
			}

			if err := s.questionRepo.Create(question); err != nil {
				return nil, err
			}

			// Create options
			for _, optReq := range qReq.Options {
				option := &domain.QuestionOption{
					QuestionID: question.ID,
					Content:    optReq.Content,
					IsCorrect:  optReq.IsCorrect,
				}
				if err := s.questionOptionRepo.Create(option); err != nil {
					return nil, err
				}
			}
		}

		// Update question count
		_ = s.quizRepo.UpdateQuestionCount(quiz.ID, len(req.Questions))
	}

	// Reload quiz with associations
	return s.GetQuizByID(quiz.ID)
}

func (s *quizService) UpdateQuiz(id, userID uint, req *request.QuizRequest) (*response.QuizResponse, error) {
	quiz, err := s.quizRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuizNotFound
		}
		return nil, err
	}

	// Check ownership
	if quiz.CreatorID != userID {
		return nil, ErrNotQuizOwner
	}

	quiz.Title = req.Title
	if req.Description != nil {
		quiz.Description = req.Description
	}
	quiz.Difficulty = req.Difficulty
	if req.IsPublic != nil {
		quiz.IsPublic = *req.IsPublic
	}

	if err := s.quizRepo.Update(quiz); err != nil {
		return nil, err
	}

	// Update categories
	if len(req.CategoryIDs) > 0 {
		// Remove old categories
		_ = s.quizRepo.RemoveCategories(quiz.ID)
		// Add new categories
		if err := s.quizRepo.AddCategories(quiz.ID, req.CategoryIDs); err != nil {
			return nil, err
		}
	}

	return s.GetQuizByID(quiz.ID)
}

func (s *quizService) DeleteQuiz(id, userID uint) error {
	quiz, err := s.quizRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuizNotFound
		}
		return err
	}

	// Check ownership
	if quiz.CreatorID != userID {
		return ErrNotQuizOwner
	}

	// Delete questions and options
	if err := s.questionRepo.DeleteByQuizID(id); err != nil {
		return err
	}

	// Decrement category quiz counts
	for _, cat := range quiz.Categories {
		_ = s.categoryRepo.DecrementQuizCount(cat.ID)
	}

	// Remove categories
	_ = s.quizRepo.RemoveCategories(id)

	return s.quizRepo.Delete(id)
}
