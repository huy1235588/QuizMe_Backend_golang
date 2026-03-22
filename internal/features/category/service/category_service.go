package service

import (
	"errors"

	"github.com/huy/quizme-backend/internal/features/category/domain"
	"github.com/huy/quizme-backend/internal/features/category/dto"
	"github.com/huy/quizme-backend/internal/features/category/repository"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound  = errors.New("category not found")
	ErrCategoryNameExists = errors.New("category name already exists")
)

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetAllCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		resp := dto.FromCategory(&cat)
		if resp != nil {
			responses[i] = *resp
		}
	}
	return responses, nil
}

func (s *categoryService) GetActiveCategories() ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}
	// Filter active categories
	var active []domain.Category
	for _, cat := range categories {
		if cat.IsActive {
			active = append(active, cat)
		}
	}
	return dto.FromCategories(active), nil
}

// GetCategoriesPaged retrieves paginated categories
func (s *categoryService) GetCategoriesPaged(page, pageSize int, search, sortBy, sortDir string) ([]dto.CategoryResponse, int64, error) {
	// For now, return all categories with no filtering
	// In a real implementation, this would use the repository with pagination
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		resp := dto.FromCategory(&cat)
		if resp != nil {
			responses[i] = *resp
		}
	}

	return responses, int64(len(responses)), nil
}

func (s *categoryService) GetCategoryByID(id uint) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return dto.FromCategory(category), nil
}

func (s *categoryService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// Check if name exists
	exists, err := s.categoryRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategoryNameExists
	}

	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		IsActive:    true, // New categories are active by default
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return dto.FromCategory(category), nil
}

func (s *categoryService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		// Check if name exists (if changed)
		if *req.Name != category.Name {
			exists, err := s.categoryRepo.ExistsByName(*req.Name)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrCategoryNameExists
			}
		}
		category.Name = *req.Name
	}

	if req.Description != nil {
		category.Description = req.Description
	}
	if req.IconURL != nil {
		category.IconURL = req.IconURL
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return dto.FromCategory(category), nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}
	return s.categoryRepo.Delete(id)
}

// IncrementPlayCount increments the play count for a category
func (s *categoryService) IncrementPlayCount(id uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	category.TotalPlayCount++
	return s.categoryRepo.Update(category)
}

// DecrementPlayCount decrements the play count for a category
func (s *categoryService) DecrementPlayCount(id uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	if category.TotalPlayCount > 0 {
		category.TotalPlayCount--
	}
	return s.categoryRepo.Update(category)
}
