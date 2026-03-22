package service

import (
	"errors"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/dto/request"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound  = errors.New("category not found")
	ErrCategoryNameExists = errors.New("category name already exists")
)

// CategoryService handles category-related operations
type CategoryService interface {
	GetAllCategories() ([]*response.CategoryResponse, error)
	GetActiveCategories() ([]*response.CategoryResponse, error)
	GetCategoryByID(id uint) (*response.CategoryResponse, error)
	CreateCategory(req *request.CategoryRequest) (*response.CategoryResponse, error)
	UpdateCategory(id uint, req *request.CategoryRequest) (*response.CategoryResponse, error)
	DeleteCategory(id uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetAllCategories() ([]*response.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return response.FromCategories(categories), nil
}

func (s *categoryService) GetActiveCategories() ([]*response.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindActive()
	if err != nil {
		return nil, err
	}
	return response.FromCategories(categories), nil
}

func (s *categoryService) GetCategoryByID(id uint) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return response.FromCategory(category), nil
}

func (s *categoryService) CreateCategory(req *request.CategoryRequest) (*response.CategoryResponse, error) {
	// Check if name exists
	exists, err := s.categoryRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategoryNameExists
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		IsActive:    isActive,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return response.FromCategory(category), nil
}

func (s *categoryService) UpdateCategory(id uint, req *request.CategoryRequest) (*response.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// Check if name exists (if changed)
	if req.Name != category.Name {
		exists, err := s.categoryRepo.ExistsByName(req.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrCategoryNameExists
		}
	}

	category.Name = req.Name
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

	return response.FromCategory(category), nil
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
