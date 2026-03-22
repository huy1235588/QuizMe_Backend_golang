package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"github.com/huy/quizme-backend/internal/dto/game"
	gameService "github.com/huy/quizme-backend/internal/features/game/service"
	gameWs "github.com/huy/quizme-backend/internal/features/game/websocket"
	roomService "github.com/huy/quizme-backend/internal/features/room/service"
	"github.com/huy/quizme-backend/internal/infra/middleware"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub                *gameWs.Hub
	authMiddleware     *middleware.AuthMiddleware
	gameSessionService gameService.GameSessionService
	roomService        roomService.RoomService
	chatService        gameService.ChatService
	participantRepo    gameService.ParticipantLookup
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(
	hub *gameWs.Hub,
	authMiddleware *middleware.AuthMiddleware,
	gameSessionService gameService.GameSessionService,
	roomService roomService.RoomService,
	chatService gameService.ChatService,
	participantRepo gameService.ParticipantLookup,
) *WebSocketHandler {
	h := &WebSocketHandler{
		hub:                hub,
		authMiddleware:     authMiddleware,
		gameSessionService: gameSessionService,
		roomService:        roomService,
		chatService:        chatService,
		participantRepo:    participantRepo,
	}

	// Set this handler as the message handler for the hub
	hub.SetMessageHandler(h)

	return h
}

// HandleConnection handles WebSocket upgrade requests
// GET /ws
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	clientID := uuid.New().String()
	client := gameWs.NewClient(h.hub, conn, clientID)

	// Try to authenticate if token provided
	token := c.Query("token")
	if token != "" {
		user, err := h.authMiddleware.ValidateToken(token)
		if err == nil && user != nil {
			client.UserID = &user.ID
		}
	}

	h.hub.Register <- client

	// Start read/write pumps
	go client.WritePump()
	go client.ReadPump()
}

// HandleMessage implements gameWs.MessageHandler
func (h *WebSocketHandler) HandleMessage(client *gameWs.Client, msg *gameWs.Message) {
	switch msg.Type {
	case gameWs.MessageTypeJoin:
		h.handleJoin(client, msg)
	case gameWs.MessageTypeLeave:
		h.handleLeave(client)
	case gameWs.MessageTypeChat:
		h.handleChat(client, msg)
	case gameWs.MessageTypeAnswer:
		h.handleAnswer(client, msg)
	case gameWs.MessageTypePing:
		h.handlePing(client)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (h *WebSocketHandler) handleJoin(client *gameWs.Client, msg *gameWs.Message) {
	var payload gameWs.JoinPayload
	if err := msg.ParsePayload(&payload); err != nil {
		h.sendError(client, "INVALID_PAYLOAD", "Invalid join payload")
		return
	}

	// Authenticate if token provided
	if payload.Token != "" && client.UserID == nil {
		user, err := h.authMiddleware.ValidateToken(payload.Token)
		if err == nil && user != nil {
			client.UserID = &user.ID
		}
	}

	// Get participant info
	var participantID uint
	var isHost bool

	if client.UserID != nil {
		participant, err := h.participantRepo.FindByRoomAndUser(payload.RoomID, *client.UserID)
		if err != nil || participant == nil {
			h.sendError(client, "NOT_IN_ROOM", "You are not in this room")
			return
		}
		participantID = participant.ID
		isHost = participant.IsHost
		client.SetRoom(payload.RoomID, client.UserID, nil, isHost)
	} else if payload.GuestName != nil {
		client.SetRoom(payload.RoomID, nil, payload.GuestName, false)
	} else {
		h.sendError(client, "AUTH_REQUIRED", "Authentication or guest name required")
		return
	}

	// Join the room
	h.hub.JoinRoom(client, payload.RoomID)

	// Send connect confirmation
	connectMsg, _ := gameWs.NewMessage(gameWs.MessageTypeConnect, map[string]interface{}{
		"clientId":      client.ID,
		"roomId":        payload.RoomID,
		"participantId": participantID,
		"isHost":        isHost,
	})
	client.SendMessage(connectMsg)

	// If game is active, send current state
	if h.gameSessionService.IsGameActive(payload.RoomID) {
		state := h.gameSessionService.GetGameState(payload.RoomID)
		stateMsg, _ := gameWs.NewMessage(gameWs.MessageTypeConnect, state)
		client.SendMessage(stateMsg)

		// Reconnect player in game session
		if participantID > 0 {
			h.gameSessionService.ReconnectPlayer(payload.RoomID, participantID, client.ID)
		}
	}

	log.Printf("Client %s joined room %d", client.ID, payload.RoomID)
}

func (h *WebSocketHandler) handleLeave(client *gameWs.Client) {
	roomID := client.GetRoomID()
	if roomID == 0 {
		return
	}

	// Get participant ID and notify game session
	if client.UserID != nil {
		participant, _ := h.participantRepo.FindByRoomAndUser(roomID, *client.UserID)
		if participant != nil {
			h.gameSessionService.DisconnectPlayer(roomID, participant.ID)
		}
	}

	h.hub.LeaveRoom(client)

	log.Printf("Client %s left room %d", client.ID, roomID)
}

func (h *WebSocketHandler) handleChat(client *gameWs.Client, msg *gameWs.Message) {
	var payload gameWs.ChatPayload
	if err := msg.ParsePayload(&payload); err != nil {
		h.sendError(client, "INVALID_PAYLOAD", "Invalid chat payload")
		return
	}

	roomID := client.GetRoomID()
	if roomID == 0 {
		h.sendError(client, "NOT_IN_ROOM", "You are not in a room")
		return
	}

	// Get sender name
	senderName := client.GetDisplayName()
	if client.UserID != nil {
		// TODO: Get username from user service
		senderName = "User"
	}

	// Broadcast chat message
	chatMsg, _ := gameWs.NewMessage(gameWs.MessageTypeChat, map[string]interface{}{
		"senderId":   client.UserID,
		"senderName": senderName,
		"content":    payload.Content,
	})
	h.hub.BroadcastToRoom(roomID, chatMsg)
}

func (h *WebSocketHandler) handleAnswer(client *gameWs.Client, msg *gameWs.Message) {
	var payload game.AnswerRequest
	if err := msg.ParsePayload(&payload); err != nil {
		h.sendError(client, "INVALID_PAYLOAD", "Invalid answer payload")
		return
	}

	roomID := client.GetRoomID()
	if roomID == 0 {
		h.sendError(client, "NOT_IN_ROOM", "You are not in a room")
		return
	}

	// Get participant ID
	if client.UserID == nil {
		h.sendError(client, "AUTH_REQUIRED", "Authentication required to submit answers")
		return
	}

	participant, err := h.participantRepo.FindByRoomAndUser(roomID, *client.UserID)
	if err != nil || participant == nil {
		h.sendError(client, "NOT_IN_ROOM", "You are not in this room")
		return
	}

	// Process answer
	err = h.gameSessionService.ProcessAnswer(roomID, participant.ID, &payload)
	if err != nil {
		h.sendError(client, "ANSWER_ERROR", err.Error())
		return
	}

	// Send acknowledgment
	ackMsg, _ := gameWs.NewMessage(gameWs.MessageTypeAnswerResult, map[string]interface{}{
		"questionId": payload.QuestionID,
		"received":   true,
	})
	client.SendMessage(ackMsg)
}

func (h *WebSocketHandler) handlePing(client *gameWs.Client) {
	pongMsg, _ := gameWs.NewMessage(gameWs.MessageTypePong, nil)
	client.SendMessage(pongMsg)
}

func (h *WebSocketHandler) sendError(client *gameWs.Client, code, message string) {
	errMsg, _ := gameWs.NewMessage(gameWs.MessageTypeError, gameWs.ErrorPayload{
		Code:    code,
		Message: message,
	})
	client.SendMessage(errMsg)
}
