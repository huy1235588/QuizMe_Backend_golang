package repository

import (
	"github.com/huy/quizme-backend/internal/domain"
	"gorm.io/gorm"
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(category *domain.Category) error
	FindByID(id uint) (*domain.Category, error)
	FindByName(name string) (*domain.Category, error)
	FindAll() ([]domain.Category, error)
	FindActive() ([]domain.Category, error)
	Update(category *domain.Category) error
	Delete(id uint) error
	ExistsByName(name string) (bool, error)
	IncrementQuizCount(id uint) error
	DecrementQuizCount(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByName(name string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindActive() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Where("is_active = ?", true).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

func (r *categoryRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Category{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *categoryRepository) IncrementQuizCount(id uint) error {
	return r.db.Model(&domain.Category{}).Where("id = ?", id).
		Update("quiz_count", gorm.Expr("quiz_count + 1")).Error
}

func (r *categoryRepository) DecrementQuizCount(id uint) error {
	return r.db.Model(&domain.Category{}).Where("id = ?", id).
		Update("quiz_count", gorm.Expr("quiz_count - 1")).Error
}
