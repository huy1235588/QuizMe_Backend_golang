package repository

import (
	"github.com/huy/quizme-backend/internal/features/category/domain"
	"gorm.io/gorm"
)

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

func (r *categoryRepository) FindAllPaged(page, pageSize int, search, sortBy, sortDir string) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var total int64

	query := r.db.Model(&domain.Category{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?)", searchPattern)
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sortBy != "" {
		order := sortBy
		if sortDir == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("name ASC")
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&categories).Error

	return categories, total, err
}

func (r *categoryRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Category{}).Count(&count).Error
	return count, err
}

func (r *categoryRepository) IncrementPlayCount(id uint) error {
	return r.db.Model(&domain.Category{}).Where("id = ?", id).
		Update("play_count", gorm.Expr("play_count + 1")).Error
}

func (r *categoryRepository) DecrementPlayCount(id uint) error {
	return r.db.Model(&domain.Category{}).Where("id = ?", id).
		Update("play_count", gorm.Expr("play_count - 1")).Error
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
