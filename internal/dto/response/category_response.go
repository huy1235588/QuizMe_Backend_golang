package response

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain"
)

// CategoryResponse represents a category in API responses
type CategoryResponse struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	IconURL        *string `json:"iconUrl,omitempty"`
	QuizCount      int     `json:"quizCount"`
	TotalPlayCount int     `json:"totalPlayCount"`
	IsActive       bool    `json:"isActive"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

// FromCategory converts a Category domain model to CategoryResponse
func FromCategory(category *domain.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:             category.ID,
		Name:           category.Name,
		Description:    category.Description,
		IconURL:        category.IconURL,
		QuizCount:      category.QuizCount,
		TotalPlayCount: category.TotalPlayCount,
		IsActive:       category.IsActive,
		CreatedAt:      category.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      category.UpdatedAt.Format(time.RFC3339),
	}
}

// FromCategories converts a slice of Category domain models to CategoryResponse slice
func FromCategories(categories []domain.Category) []*CategoryResponse {
	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = FromCategory(&category)
	}
	return responses
}
