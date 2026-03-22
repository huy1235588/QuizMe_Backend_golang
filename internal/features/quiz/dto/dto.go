package dto

import (
	"time"

	"github.com/huy/quizme-backend/internal/features/quiz/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
)

// ==== REQUEST DTOs ====

// QuizRequest represents the request body for quiz operations
type QuizRequest struct {
	Title       string                `json:"title" validate:"required,max=100"`
	Description *string               `json:"description" validate:"max=1000"`
	CategoryIDs []uint                `json:"categoryIds"`
	Difficulty  enums.Difficulty      `json:"difficulty" validate:"required"`
	IsPublic    *bool                 `json:"isPublic"`
	Questions   []QuestionRequest     `json:"questions,omitempty"`
}

// QuestionRequest represents the request body for question operations
type QuestionRequest struct {
	Content     string                   `json:"content" validate:"required"`
	ImageURL    *string                  `json:"imageUrl"`
	VideoURL    *string                  `json:"videoUrl"`
	AudioURL    *string                  `json:"audioUrl"`
	FunFact     *string                  `json:"funFact"`
	Explanation *string                  `json:"explanation"`
	TimeLimit   int                      `json:"timeLimit" validate:"min=5,max=300"`
	Points      int                      `json:"points" validate:"min=1,max=1000"`
	OrderNumber int                      `json:"orderNumber" validate:"min=1"`
	Type        enums.QuestionType       `json:"type" validate:"required"`
	Options     []QuestionOptionRequest  `json:"options" validate:"required,min=2,dive"`
}

// QuestionOptionRequest represents the request body for question option operations
type QuestionOptionRequest struct {
	Content   string `json:"content" validate:"required"`
	IsCorrect bool   `json:"isCorrect"`
}

// BatchQuestionRequest represents the request body for batch question creation
type BatchQuestionRequest struct {
	QuizID    uint              `json:"quizId" validate:"required"`
	Questions []QuestionRequest `json:"questions" validate:"required,min=1,dive"`
}

// ==== RESPONSE DTOs ====

// QuizResponse represents a quiz in API responses
type QuizResponse struct {
	ID             uint             `json:"id"`
	Title          string           `json:"title"`
	Description    *string          `json:"description,omitempty"`
	QuizThumbnails *string          `json:"quizThumbnails,omitempty"`
	CategoryIDs    []uint           `json:"categoryIds,omitempty"`
	CategoryNames  []string         `json:"categoryNames,omitempty"`
	CreatorID      uint             `json:"creatorId"`
	CreatorName    string           `json:"creatorName"`
	CreatorAvatar  *string          `json:"creatorAvatar,omitempty"`
	Difficulty     enums.Difficulty `json:"difficulty"`
	IsPublic       bool             `json:"isPublic"`
	PlayCount      int              `json:"playCount"`
	QuestionCount  int              `json:"questionCount"`
	FavoriteCount  int              `json:"favoriteCount"`
	CreatedAt      string           `json:"createdAt"`
	UpdatedAt      string           `json:"updatedAt"`
}

// QuestionResponse represents a question in API responses
type QuestionResponse struct {
	ID          uint                       `json:"id"`
	QuizID      uint                       `json:"quizId"`
	Content     string                     `json:"content"`
	ImageURL    *string                    `json:"imageUrl,omitempty"`
	VideoURL    *string                    `json:"videoUrl,omitempty"`
	AudioURL    *string                    `json:"audioUrl,omitempty"`
	FunFact     *string                    `json:"funFact,omitempty"`
	Explanation *string                    `json:"explanation,omitempty"`
	TimeLimit   int                        `json:"timeLimit"`
	Points      int                        `json:"points"`
	OrderNumber int                        `json:"orderNumber"`
	Type        enums.QuestionType         `json:"type"`
	Options     []*QuestionOptionResponse  `json:"options,omitempty"`
	CreatedAt   string                     `json:"createdAt"`
	UpdatedAt   string                     `json:"updatedAt"`
}

// QuestionOptionResponse represents a question option in API responses
type QuestionOptionResponse struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"isCorrect"`
}

// ==== CONVERSION FUNCTIONS ====

// FromQuiz converts a Quiz domain model to QuizResponse
func FromQuiz(quiz *domain.Quiz) *QuizResponse {
	response := &QuizResponse{
		ID:             quiz.ID,
		Title:          quiz.Title,
		Description:    quiz.Description,
		QuizThumbnails: quiz.QuizThumbnails,
		CreatorID:      quiz.CreatorID,
		Difficulty:     quiz.Difficulty,
		IsPublic:       quiz.IsPublic,
		PlayCount:      quiz.PlayCount,
		QuestionCount:  quiz.QuestionCount,
		FavoriteCount:  quiz.FavoriteCount,
		CreatedAt:      quiz.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      quiz.UpdatedAt.Format(time.RFC3339),
	}

	// Add creator info if available
	if quiz.Creator != nil {
		response.CreatorName = quiz.Creator.FullName
		response.CreatorAvatar = quiz.Creator.ProfileImage
	}

	// Add category info if available
	if len(quiz.Categories) > 0 {
		response.CategoryIDs = make([]uint, len(quiz.Categories))
		response.CategoryNames = make([]string, len(quiz.Categories))
		for i, cat := range quiz.Categories {
			response.CategoryIDs[i] = cat.ID
			response.CategoryNames[i] = cat.Name
		}
	}

	return response
}

// FromQuizzes converts a slice of Quiz domain models to QuizResponse slice
func FromQuizzes(quizzes []domain.Quiz) []*QuizResponse {
	responses := make([]*QuizResponse, len(quizzes))
	for i, quiz := range quizzes {
		responses[i] = FromQuiz(&quiz)
	}
	return responses
}

// FromQuestion converts a Question domain model to QuestionResponse
func FromQuestion(question *domain.Question) *QuestionResponse {
	response := &QuestionResponse{
		ID:          question.ID,
		QuizID:      question.QuizID,
		Content:     question.Content,
		ImageURL:    question.ImageURL,
		VideoURL:    question.VideoURL,
		AudioURL:    question.AudioURL,
		FunFact:     question.FunFact,
		Explanation: question.Explanation,
		TimeLimit:   question.TimeLimit,
		Points:      question.Points,
		OrderNumber: question.OrderNumber,
		Type:        question.Type,
		CreatedAt:   question.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   question.UpdatedAt.Format(time.RFC3339),
	}

	// Add options if available
	if len(question.Options) > 0 {
		response.Options = make([]*QuestionOptionResponse, len(question.Options))
		for i, opt := range question.Options {
			response.Options[i] = FromQuestionOption(&opt)
		}
	}

	return response
}

// FromQuestionOption converts a QuestionOption domain model to QuestionOptionResponse
func FromQuestionOption(option *domain.QuestionOption) *QuestionOptionResponse {
	return &QuestionOptionResponse{
		ID:        option.ID,
		Content:   option.Content,
		IsCorrect: option.IsCorrect,
	}
}

// FromQuestions converts a slice of Question domain models to QuestionResponse slice
func FromQuestions(questions []domain.Question) []*QuestionResponse {
	responses := make([]*QuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = FromQuestion(&question)
	}
	return responses
}
