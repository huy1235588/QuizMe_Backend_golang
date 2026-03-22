package request

import "github.com/huy/quizme-backend/internal/domain/enums"

// QuizRequest represents the request body for quiz operations
type QuizRequest struct {
	Title       string           `json:"title" validate:"required,max=100"`
	Description *string          `json:"description" validate:"max=1000"`
	CategoryIDs []uint           `json:"categoryIds"`
	Difficulty  enums.Difficulty `json:"difficulty" validate:"required"`
	IsPublic    *bool            `json:"isPublic"`
	Questions   []QuestionRequest `json:"questions,omitempty"`
}

// QuestionRequest represents the request body for question operations
type QuestionRequest struct {
	Content     string                  `json:"content" validate:"required"`
	ImageURL    *string                 `json:"imageUrl"`
	VideoURL    *string                 `json:"videoUrl"`
	AudioURL    *string                 `json:"audioUrl"`
	FunFact     *string                 `json:"funFact"`
	Explanation *string                 `json:"explanation"`
	TimeLimit   int                     `json:"timeLimit" validate:"min=5,max=300"`
	Points      int                     `json:"points" validate:"min=1,max=1000"`
	OrderNumber int                     `json:"orderNumber" validate:"min=1"`
	Type        enums.QuestionType      `json:"type" validate:"required"`
	Options     []QuestionOptionRequest `json:"options" validate:"required,min=2,dive"`
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
