package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huy/quizme-backend/internal/dto/response"
	gameService "github.com/huy/quizme-backend/internal/features/game/service"
	roomService "github.com/huy/quizme-backend/internal/features/room/service"
	"github.com/huy/quizme-backend/internal/infra/middleware"
)

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	gameSessionService gameService.GameSessionService
	roomService        roomService.RoomService
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameSessionService gameService.GameSessionService, roomService roomService.RoomService) *GameHandler {
	return &GameHandler{
		gameSessionService: gameSessionService,
		roomService:        roomService,
	}
}

// GetGameState handles getting the current game state
// GET /api/game/state/:roomId
func (h *GameHandler) GetGameState(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid room ID"))
		return
	}

	state := h.gameSessionService.GetGameState(uint(roomID))
	c.JSON(http.StatusOK, response.Success(state, "Game state retrieved successfully"))
}

// InitGame handles initializing a game session
// POST /api/game/init/:roomId
func (h *GameHandler) InitGame(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
		return
	}

	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid room ID"))
		return
	}

	// Verify user is the host
	room, err := h.roomService.GetRoomByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Room not found"))
		return
	}

	if room.HostID != currentUser.ID {
		c.JSON(http.StatusForbidden, response.Error("Only the host can initialize the game"))
		return
	}

	session, err := h.gameSessionService.InitGameSession(uint(roomID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to initialize game session"))
		return
	}

	c.JSON(http.StatusOK, response.Success(map[string]interface{}{
		"roomId":           session.RoomID,
		"quizId":           session.QuizID,
		"totalQuestions":   len(session.Questions),
		"participantCount": len(session.Participants),
	}, "Game session initialized successfully"))
}

// StartGame handles starting a game
// POST /api/game/start/:roomId
func (h *GameHandler) StartGame(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
		return
	}

	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid room ID"))
		return
	}

	// Verify user is the host
	room, err := h.roomService.GetRoomByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Room not found"))
		return
	}

	if room.HostID != currentUser.ID {
		c.JSON(http.StatusForbidden, response.Error("Only the host can start the game"))
		return
	}

	// Check if game session exists, if not initialize it
	if h.gameSessionService.GetSession(uint(roomID)) == nil {
		_, err := h.gameSessionService.InitGameSession(uint(roomID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error("Failed to initialize game session"))
			return
		}
	}

	if err := h.gameSessionService.StartGame(uint(roomID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to start game"))
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "Game started successfully"))
}
