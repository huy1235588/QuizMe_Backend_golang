package response

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
)

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
