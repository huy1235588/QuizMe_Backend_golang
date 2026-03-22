package service

import (
	"github.com/huy/quizme-backend/internal/features/category/dto"
)

// CategoryService defines the interface for category operations
type CategoryService interface {
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategoryByID(id uint) (*dto.CategoryResponse, error)
	GetAllCategories() ([]dto.CategoryResponse, error)
	GetActiveCategories() ([]*dto.CategoryResponse, error)
	GetCategoriesPaged(page, pageSize int, search, sortBy, sortDir string) ([]dto.CategoryResponse, int64, error)
	UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id uint) error
	IncrementPlayCount(id uint) error
	DecrementPlayCount(id uint) error
}
