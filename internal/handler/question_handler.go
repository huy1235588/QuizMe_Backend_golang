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

// QuestionHandler handles question-related HTTP requests
type QuestionHandler struct {
	questionService service.QuestionService
}

// NewQuestionHandler creates a new question handler
func NewQuestionHandler(questionService service.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

// GetQuestionByID handles getting a question by ID
// GET /api/questions/:id
func (h *QuestionHandler) GetQuestionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid question ID"))
		return
	}

	question, err := h.questionService.GetQuestionByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Question not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get question"))
		return
	}

	c.JSON(http.StatusOK, response.Success(question, "Question retrieved successfully"))
}

// GetQuestionsByQuizID handles getting questions by quiz ID
// GET /api/questions/quiz/:quizId
func (h *QuestionHandler) GetQuestionsByQuizID(c *gin.Context) {
	quizID, err := strconv.ParseUint(c.Param("quizId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid quiz ID"))
		return
	}

	questions, err := h.questionService.GetQuestionsByQuizID(uint(quizID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get questions"))
		return
	}

	c.JSON(http.StatusOK, response.Success(questions, "Questions retrieved successfully"))
}

// CreateQuestion handles question creation
// POST /api/questions
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req struct {
		QuizID uint `json:"quizId" binding:"required"`
		request.QuestionRequest
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	question, err := h.questionService.CreateQuestion(&req.QuestionRequest, req.QuizID)
	if err != nil {
		if errors.Is(err, service.ErrQuizNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Quiz not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to create question"))
		return
	}

	c.JSON(http.StatusCreated, response.Success(question, "Question created successfully"))
}

// CreateBatchQuestions handles batch question creation
// POST /api/questions/batch
func (h *QuestionHandler) CreateBatchQuestions(c *gin.Context) {
	var req request.BatchQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	questions, err := h.questionService.CreateBatchQuestions(&req)
	if err != nil {
		if errors.Is(err, service.ErrQuizNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Quiz not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to create questions"))
		return
	}

	c.JSON(http.StatusCreated, response.Success(questions, "Questions created successfully"))
}

// UpdateQuestion handles question update
// PUT /api/questions/:id
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid question ID"))
		return
	}

	var req request.QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	question, err := h.questionService.UpdateQuestion(uint(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Question not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to update question"))
		return
	}

	c.JSON(http.StatusOK, response.Success(question, "Question updated successfully"))
}

// DeleteQuestion handles question deletion
// DELETE /api/questions/:id
func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid question ID"))
		return
	}

	if err := h.questionService.DeleteQuestion(uint(id)); err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Question not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to delete question"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "Question deleted successfully"))
}
