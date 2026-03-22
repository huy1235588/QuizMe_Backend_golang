package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huy/quizme-backend/internal/dto/request"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/service"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetAllCategories handles getting all categories
// GET /api/categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get categories"))
		return
	}
	c.JSON(http.StatusOK, response.Success(categories, "Categories retrieved successfully"))
}

// GetActiveCategories handles getting active categories
// GET /api/categories/active
func (h *CategoryHandler) GetActiveCategories(c *gin.Context) {
	categories, err := h.categoryService.GetActiveCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get categories"))
		return
	}
	c.JSON(http.StatusOK, response.Success(categories, "Active categories retrieved successfully"))
}

// GetCategoryByID handles getting a category by ID
// GET /api/categories/:id
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid category ID"))
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Category not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get category"))
		return
	}

	c.JSON(http.StatusOK, response.Success(category, "Category retrieved successfully"))
}

// CreateCategory handles category creation (admin only)
// POST /api/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req request.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNameExists) {
			c.JSON(http.StatusConflict, response.Error("Category name already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to create category"))
		return
	}

	c.JSON(http.StatusCreated, response.Success(category, "Category created successfully"))
}

// UpdateCategory handles category update (admin only)
// PUT /api/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid category ID"))
		return
	}

	var req request.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	category, err := h.categoryService.UpdateCategory(uint(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Category not found"))
			return
		}
		if errors.Is(err, service.ErrCategoryNameExists) {
			c.JSON(http.StatusConflict, response.Error("Category name already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to update category"))
		return
	}

	c.JSON(http.StatusOK, response.Success(category, "Category updated successfully"))
}

// DeleteCategory handles category deletion (admin only)
// DELETE /api/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid category ID"))
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Category not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to delete category"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "Category deleted successfully"))
}
