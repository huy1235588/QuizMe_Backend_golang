package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	roomdto "github.com/huy/quizme-backend/internal/features/room/dto"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/infra/middleware"
	roomservice "github.com/huy/quizme-backend/internal/features/room/service"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	chatService roomservice.ChatService
}

// NewChatHandler creates a new chat handler
func NewChatHandler(chatService roomservice.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// GetChatHistory handles getting chat history for a room
// GET /api/chat/room/:roomId
func (h *ChatHandler) GetChatHistory(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid room ID"))
		return
	}

	limit := 100 // default
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	messages, err := h.chatService.GetChatHistory(uint(roomID), limit)
	if err != nil {
		if errors.Is(err, roomservice.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Room not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get chat history"))
		return
	}

	c.JSON(http.StatusOK, response.Success(messages, "Chat history retrieved successfully"))
}

// SendMessage handles sending a chat message
// POST /api/chat/send
func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req roomdto.ChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	var userID *uint
	currentUser := middleware.GetCurrentUser(c)
	if currentUser != nil {
		userID = &currentUser.ID
	}

	message, err := h.chatService.SendMessage(&req, userID)
	if err != nil {
		if errors.Is(err, roomservice.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Room not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to send message"))
		return
	}

	c.JSON(http.StatusCreated, response.Success(message, "Message sent successfully"))
}
