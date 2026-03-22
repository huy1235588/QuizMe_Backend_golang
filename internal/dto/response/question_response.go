package response

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
)

// QuestionResponse represents a question in API responses
type QuestionResponse struct {
	ID          uint                     `json:"id"`
	QuizID      uint                     `json:"quizId"`
	Content     string                   `json:"content"`
	ImageURL    *string                  `json:"imageUrl,omitempty"`
	VideoURL    *string                  `json:"videoUrl,omitempty"`
	AudioURL    *string                  `json:"audioUrl,omitempty"`
	FunFact     *string                  `json:"funFact,omitempty"`
	Explanation *string                  `json:"explanation,omitempty"`
	TimeLimit   int                      `json:"timeLimit"`
	Points      int                      `json:"points"`
	OrderNumber int                      `json:"orderNumber"`
	Type        enums.QuestionType       `json:"type"`
	Options     []*QuestionOptionResponse `json:"options,omitempty"`
	CreatedAt   string                   `json:"createdAt"`
	UpdatedAt   string                   `json:"updatedAt"`
}

// QuestionOptionResponse represents a question option in API responses
type QuestionOptionResponse struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"isCorrect"`
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
