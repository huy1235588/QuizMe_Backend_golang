package repository

import (
	"github.com/huy/quizme-backend/internal/features/category/domain"
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(category *domain.Category) error
	FindByID(id uint) (*domain.Category, error)
	FindByName(name string) (*domain.Category, error)
	FindAll() ([]domain.Category, error)
	FindAllPaged(page, pageSize int, search, sortBy, sortDir string) ([]domain.Category, int64, error)
	Update(category *domain.Category) error
	Delete(id uint) error
	Count() (int64, error)
	ExistsByName(name string) (bool, error)
	IncrementPlayCount(id uint) error
	DecrementPlayCount(id uint) error
	IncrementQuizCount(id uint) error
	DecrementQuizCount(id uint) error
}
