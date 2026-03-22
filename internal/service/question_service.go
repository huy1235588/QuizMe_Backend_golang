package service

import (
	"errors"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/dto/request"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/repository"
	"gorm.io/gorm"
)

// QuestionService handles question-related operations
type QuestionService interface {
	GetAllQuestions() ([]*response.QuestionResponse, error)
	GetQuestionByID(id uint) (*response.QuestionResponse, error)
	GetQuestionsByQuizID(quizID uint) ([]*response.QuestionResponse, error)
	CreateQuestion(req *request.QuestionRequest, quizID uint) (*response.QuestionResponse, error)
	CreateBatchQuestions(req *request.BatchQuestionRequest) ([]*response.QuestionResponse, error)
	UpdateQuestion(id uint, req *request.QuestionRequest) (*response.QuestionResponse, error)
	DeleteQuestion(id uint) error
}

type questionService struct {
	questionRepo       repository.QuestionRepository
	questionOptionRepo repository.QuestionOptionRepository
	quizRepo           repository.QuizRepository
}

// NewQuestionService creates a new question service
func NewQuestionService(
	questionRepo repository.QuestionRepository,
	questionOptionRepo repository.QuestionOptionRepository,
	quizRepo repository.QuizRepository,
) QuestionService {
	return &questionService{
		questionRepo:       questionRepo,
		questionOptionRepo: questionOptionRepo,
		quizRepo:           quizRepo,
	}
}

func (s *questionService) GetAllQuestions() ([]*response.QuestionResponse, error) {
	// This is typically not used, but included for completeness
	return nil, errors.New("not implemented")
}

func (s *questionService) GetQuestionByID(id uint) (*response.QuestionResponse, error) {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}
	return response.FromQuestion(question), nil
}

func (s *questionService) GetQuestionsByQuizID(quizID uint) ([]*response.QuestionResponse, error) {
	questions, err := s.questionRepo.FindByQuizID(quizID)
	if err != nil {
		return nil, err
	}
	return response.FromQuestions(questions), nil
}

func (s *questionService) CreateQuestion(req *request.QuestionRequest, quizID uint) (*response.QuestionResponse, error) {
	// Verify quiz exists
	_, err := s.quizRepo.FindByID(quizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuizNotFound
		}
		return nil, err
	}

	question := &domain.Question{
		QuizID:      quizID,
		Content:     req.Content,
		ImageURL:    req.ImageURL,
		VideoURL:    req.VideoURL,
		AudioURL:    req.AudioURL,
		FunFact:     req.FunFact,
		Explanation: req.Explanation,
		TimeLimit:   req.TimeLimit,
		Points:      req.Points,
		OrderNumber: req.OrderNumber,
		Type:        req.Type,
	}

	if err := s.questionRepo.Create(question); err != nil {
		return nil, err
	}

	// Create options
	for _, optReq := range req.Options {
		option := &domain.QuestionOption{
			QuestionID: question.ID,
			Content:    optReq.Content,
			IsCorrect:  optReq.IsCorrect,
		}
		if err := s.questionOptionRepo.Create(option); err != nil {
			return nil, err
		}
	}

	// Update quiz question count
	count, _ := s.questionRepo.CountByQuizID(quizID)
	_ = s.quizRepo.UpdateQuestionCount(quizID, int(count))

	// Reload question with options
	return s.GetQuestionByID(question.ID)
}

func (s *questionService) CreateBatchQuestions(req *request.BatchQuestionRequest) ([]*response.QuestionResponse, error) {
	// Verify quiz exists
	_, err := s.quizRepo.FindByID(req.QuizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuizNotFound
		}
		return nil, err
	}

	responses := make([]*response.QuestionResponse, len(req.Questions))

	for i, qReq := range req.Questions {
		question := &domain.Question{
			QuizID:      req.QuizID,
			Content:     qReq.Content,
			ImageURL:    qReq.ImageURL,
			VideoURL:    qReq.VideoURL,
			AudioURL:    qReq.AudioURL,
			FunFact:     qReq.FunFact,
			Explanation: qReq.Explanation,
			TimeLimit:   qReq.TimeLimit,
			Points:      qReq.Points,
			OrderNumber: qReq.OrderNumber,
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

		// Reload with options
		resp, _ := s.GetQuestionByID(question.ID)
		responses[i] = resp
	}

	// Update quiz question count
	count, _ := s.questionRepo.CountByQuizID(req.QuizID)
	_ = s.quizRepo.UpdateQuestionCount(req.QuizID, int(count))

	return responses, nil
}

func (s *questionService) UpdateQuestion(id uint, req *request.QuestionRequest) (*response.QuestionResponse, error) {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}

	question.Content = req.Content
	question.ImageURL = req.ImageURL
	question.VideoURL = req.VideoURL
	question.AudioURL = req.AudioURL
	question.FunFact = req.FunFact
	question.Explanation = req.Explanation
	question.TimeLimit = req.TimeLimit
	question.Points = req.Points
	question.OrderNumber = req.OrderNumber
	question.Type = req.Type

	if err := s.questionRepo.Update(question); err != nil {
		return nil, err
	}

	// Update options - delete old and create new
	if err := s.questionOptionRepo.DeleteByQuestionID(id); err != nil {
		return nil, err
	}

	for _, optReq := range req.Options {
		option := &domain.QuestionOption{
			QuestionID: question.ID,
			Content:    optReq.Content,
			IsCorrect:  optReq.IsCorrect,
		}
		if err := s.questionOptionRepo.Create(option); err != nil {
			return nil, err
		}
	}

	return s.GetQuestionByID(id)
}

func (s *questionService) DeleteQuestion(id uint) error {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotFound
		}
		return err
	}

	quizID := question.QuizID

	if err := s.questionRepo.Delete(id); err != nil {
		return err
	}

	// Update quiz question count
	count, _ := s.questionRepo.CountByQuizID(quizID)
	_ = s.quizRepo.UpdateQuestionCount(quizID, int(count))

	return nil
}
