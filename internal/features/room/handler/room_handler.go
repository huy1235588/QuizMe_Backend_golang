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

// RoomHandler handles room-related HTTP requests
type RoomHandler struct {
	roomService roomservice.RoomService
}

// NewRoomHandler creates a new room handler
func NewRoomHandler(roomService roomservice.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

// CreateRoom handles room creation
// POST /api/rooms
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, response.Error("Authentication required"))
		return
	}

	var req roomdto.RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	room, err := h.roomService.CreateRoom(currentUser.ID, &req)
	if err != nil {
		if errors.Is(err, roomservice.ErrQuizNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Quiz not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to create room"))
		return
	}

	c.JSON(http.StatusCreated, response.Success(room, "Room created successfully"))
}

// GetRoomByCode handles getting a room by code
// GET /api/rooms/:code
func (h *RoomHandler) GetRoomByCode(c *gin.Context) {
	code := c.Param("code")

	room, err := h.roomService.GetRoomByCode(code)
	if err != nil {
		if errors.Is(err, roomservice.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Room not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get room"))
		return
	}

	c.JSON(http.StatusOK, response.Success(room, "Room retrieved successfully"))
}

// GetWaitingRooms handles getting waiting rooms
// GET /api/rooms/waiting
func (h *RoomHandler) GetWaitingRooms(c *gin.Context) {
	rooms, err := h.roomService.GetWaitingRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get rooms"))
		return
	}

	c.JSON(http.StatusOK, response.Success(rooms, "Rooms retrieved successfully"))
}

// GetAvailableRooms handles getting available rooms with filters
// GET /api/rooms/available
func (h *RoomHandler) GetAvailableRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var search *string
	if s := c.Query("search"); s != "" {
		search = &s
	}

	var quizID *uint
	if qIDStr := c.Query("quizId"); qIDStr != "" {
		if qID, err := strconv.ParseUint(qIDStr, 10, 32); err == nil {
			id := uint(qID)
			quizID = &id
		}
	}

	rooms, total, err := h.roomService.GetAvailableRooms(search, quizID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to get rooms"))
		return
	}

	pageResponse := response.PageResponse[[]*roomdto.RoomResponse]{
		Content:       rooms,
		Page:          page,
		Size:          pageSize,
		TotalElements: total,
		TotalPages:    (total + int64(pageSize) - 1) / int64(pageSize),
	}

	c.JSON(http.StatusOK, response.Success(pageResponse, "Available rooms retrieved successfully"))
}

// JoinRoomByCode handles joining a room by code
// POST /api/rooms/join
func (h *RoomHandler) JoinRoomByCode(c *gin.Context) {
	var req roomdto.JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	var userID *uint
	currentUser := middleware.GetCurrentUser(c)
	if currentUser != nil {
		userID = &currentUser.ID
	}

	room, participant, err := h.roomService.JoinRoomByCode(req.Code, userID, req.GuestName, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, roomservice.ErrRoomNotFound):
			c.JSON(http.StatusNotFound, response.Error("Room not found"))
		case errors.Is(err, roomservice.ErrRoomFull):
			c.JSON(http.StatusBadRequest, response.Error("Room is full"))
		case errors.Is(err, roomservice.ErrWrongPassword):
			c.JSON(http.StatusUnauthorized, response.Error("Wrong room password"))
		case errors.Is(err, roomservice.ErrRoomNotWaiting):
			c.JSON(http.StatusBadRequest, response.Error("Room is not accepting new players"))
		case errors.Is(err, roomservice.ErrAlreadyInRoom):
			c.JSON(http.StatusBadRequest, response.Error("You are already in this room"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("Failed to join room"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{"room": room, "participant": participant}, "Joined room successfully"))
}

// JoinRoomByID handles joining a room by ID
// POST /api/rooms/join/:roomId
func (h *RoomHandler) JoinRoomByID(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid room ID"))
		return
	}

	var req roomdto.JoinRoomByIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body for authenticated users
		req = roomdto.JoinRoomByIDRequest{}
	}

	var userID *uint
	currentUser := middleware.GetCurrentUser(c)
	if currentUser != nil {
		userID = &currentUser.ID
	}

	participant, err := h.roomService.JoinRoom(uint(roomID), userID, req.GuestName, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, roomservice.ErrRoomNotFound):
			c.JSON(http.StatusNotFound, response.Error("Room not found"))
		case errors.Is(err, roomservice.ErrRoomFull):
			c.JSON(http.StatusBadRequest, response.Error("Room is full"))
		case errors.Is(err, roomservice.ErrWrongPassword):
			c.JSON(http.StatusUnauthorized, response.Error("Wrong room password"))
		case errors.Is(err, roomservice.ErrRoomNotWaiting):
			c.JSON(http.StatusBadRequest, response.Error("Room is not accepting new players"))
		case errors.Is(err, roomservice.ErrAlreadyInRoom):
			c.JSON(http.StatusBadRequest, response.Error("You are already in this room"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("Failed to join room"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(participant, "Joined room successfully"))
}

// LeaveRoom handles leaving a room
// DELETE /api/rooms/leave/:roomId
func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid room ID"))
		return
	}

	var userID *uint
	currentUser := middleware.GetCurrentUser(c)
	if currentUser != nil {
		userID = &currentUser.ID
	}

	if err := h.roomService.LeaveRoom(uint(roomID), userID, nil); err != nil {
		switch {
		case errors.Is(err, roomservice.ErrRoomNotFound):
			c.JSON(http.StatusNotFound, response.Error("Room not found"))
		case errors.Is(err, roomservice.ErrNotInRoom):
			c.JSON(http.StatusBadRequest, response.Error("You are not in this room"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("Failed to leave room"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "Left room successfully"))
}

// CloseRoom handles closing a room (host only)
// PATCH /api/rooms/close/:roomId
func (h *RoomHandler) CloseRoom(c *gin.Context) {
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

	if err := h.roomService.CloseRoom(uint(roomID), currentUser.ID); err != nil {
		switch {
		case errors.Is(err, roomservice.ErrRoomNotFound):
			c.JSON(http.StatusNotFound, response.Error("Room not found"))
		case errors.Is(err, roomservice.ErrNotRoomHost):
			c.JSON(http.StatusForbidden, response.Error("You are not the host of this room"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("Failed to close room"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success[any](nil, "Room closed successfully"))
}

// UpdateRoom handles updating a room (host only)
// PATCH /api/rooms/:roomId
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
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

	var req roomdto.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body"))
		return
	}

	room, err := h.roomService.UpdateRoom(uint(roomID), currentUser.ID, &req)
	if err != nil {
		switch {
		case errors.Is(err, roomservice.ErrRoomNotFound):
			c.JSON(http.StatusNotFound, response.Error("Room not found"))
		case errors.Is(err, roomservice.ErrNotRoomHost):
			c.JSON(http.StatusForbidden, response.Error("You are not the host of this room"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("Failed to update room"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(room, "Room updated successfully"))
}

// StartGame handles starting a game (host only)
// POST /api/rooms/start/:roomId
func (h *RoomHandler) StartGame(c *gin.Context) {
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

	room, err := h.roomService.StartGame(uint(roomID), currentUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, roomservice.ErrRoomNotFound):
			c.JSON(http.StatusNotFound, response.Error("Room not found"))
		case errors.Is(err, roomservice.ErrNotRoomHost):
			c.JSON(http.StatusForbidden, response.Error("You are not the host of this room"))
		case errors.Is(err, roomservice.ErrRoomNotWaiting):
			c.JSON(http.StatusBadRequest, response.Error("Game has already started"))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("Failed to start game"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(room, "Game started successfully"))
}
