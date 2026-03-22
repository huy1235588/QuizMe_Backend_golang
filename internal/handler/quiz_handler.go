package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/dto/request"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/middleware"
	"github.com/huy/quizme-backend/internal/service"
)

// QuizHandler handles quiz-related HTTP requests
type QuizHandler struct {
	quizService service.QuizService
}

// NewQuizHandler creates a new quiz handler
func NewQuizHandler(quizService service.QuizService) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
	}
}

// GetAllQuizzes handles getting all quizzes
// GET /api/quizzes
func (h *QuizHandler) GetAllQuizzes(c *gin.Context) {
	quizzes, err := h.quizService.GetAllQuizzes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get quizzes"))
		return
	}
	c.JSON(http.StatusOK, response.Success(quizzes, "Quizzes retrieved successfully"))
}

// GetPublicQuizzes handles getting public quizzes
// GET /api/quizzes/public
func (h *QuizHandler) GetPublicQuizzes(c *gin.Context) {
	quizzes, err := h.quizService.GetPublicQuizzes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get quizzes"))
		return
	}
	c.JSON(http.StatusOK, response.Success(quizzes, "Public quizzes retrieved successfully"))
}

// GetQuizByID handles getting a quiz by ID
// GET /api/quizzes/:id
func (h *QuizHandler) GetQuizByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid quiz ID"))
		return
	}

	quiz, err := h.quizService.GetQuizByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrQuizNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Quiz not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get quiz"))
		return
	}

	c.JSON(http.StatusOK, response.Success(quiz, "Quiz retrieved successfully"))
}

// GetQuizzesByDifficulty handles getting quizzes by difficulty
// GET /api/quizzes/difficulty/:difficulty
func (h *QuizHandler) GetQuizzesByDifficulty(c *gin.Context) {
	difficultyStr := c.Param("difficulty")
	difficulty := enums.Difficulty(difficultyStr)

	if !difficulty.IsValid() {
		c.JSON(http.StatusBadRequest, response.Error("Invalid difficulty level"))
		return
	}

	quizzes, err := h.quizService.GetQuizzesByDifficulty(difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get quizzes"))
		return
	}

	c.JSON(http.StatusOK, response.Success(quizzes, "Quizzes retrieved successfully"))
}

// GetPagedQuizzes handles getting paginated quizzes with filters
// GET /api/quizzes/paged
func (h *QuizHandler) GetPagedQuizzes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var categoryID *uint
	if catIDStr := c.Query("categoryId"); catIDStr != "" {
		if catID, err := strconv.ParseUint(catIDStr, 10, 32); err == nil {
			id := uint(catID)
			categoryID = &id
		}
	}

	var difficulty *string
	if d := c.Query("difficulty"); d != "" {
		difficulty = &d
	}

	var isPublic *bool
	if p := c.Query("isPublic"); p != "" {
		val := p == "true"
		isPublic = &val
	}

	var search *string
	if s := c.Query("search"); s != "" {
		search = &s
	}

	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortDir := c.DefaultQuery("sortDir", "desc")

	quizzes, total, err := h.quizService.GetQuizzesWithFilters(categoryID, difficulty, isPublic, search, page, pageSize, sortBy, sortDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get quizzes"))
		return
	}

	pageResponse := response.PageResponse[*response.QuizResponse]{
		Content:       quizzes,
		Page:          page,
		Size:          pageSize,
		TotalElements: total,
		TotalPages:    (total + int64(pageSize) - 1) / int64(pageSize),
	}

	c.JSON(http.StatusOK, response.Success(pageResponse, "Quizzes retrieved successfully"))
}

// CreateQuiz handles quiz creation
// POST /api/quizzes
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
		return
	}

	var req request.QuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	quiz, err := h.quizService.CreateQuiz(currentUser.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to create quiz"))
		return
	}

	c.JSON(http.StatusCreated, response.Success(quiz, "Quiz created successfully"))
}

// UpdateQuiz handles quiz update
// PUT /api/quizzes/:id
func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid quiz ID"))
		return
	}

	var req request.QuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	quiz, err := h.quizService.UpdateQuiz(uint(id), currentUser.ID, &req)
	if err != nil {
		if errors.Is(err, service.ErrQuizNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Quiz not found"))
			return
		}
		if errors.Is(err, service.ErrNotQuizOwner) {
			c.JSON(http.StatusForbidden, response.Error("You are not the owner of this quiz"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to update quiz"))
		return
	}

	c.JSON(http.StatusOK, response.Success(quiz, "Quiz updated successfully"))
}

// DeleteQuiz handles quiz deletion
// DELETE /api/quizzes/:id
func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid quiz ID"))
		return
	}

	if err := h.quizService.DeleteQuiz(uint(id), currentUser.ID); err != nil {
		if errors.Is(err, service.ErrQuizNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Quiz not found"))
			return
		}
		if errors.Is(err, service.ErrNotQuizOwner) {
			c.JSON(http.StatusForbidden, response.Error("You are not the owner of this quiz"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to delete quiz"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "Quiz deleted successfully"))
}
