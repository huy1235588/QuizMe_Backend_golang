package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"github.com/huy/quizme-backend/internal/dto/game"
	"github.com/huy/quizme-backend/internal/middleware"
	"github.com/huy/quizme-backend/internal/service"
	"github.com/huy/quizme-backend/internal/websocket"
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
	hub                *websocket.Hub
	authMiddleware     *middleware.AuthMiddleware
	gameSessionService service.GameSessionService
	roomService        service.RoomService
	chatService        service.ChatService
	participantRepo    service.ParticipantLookup
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(
	hub *websocket.Hub,
	authMiddleware *middleware.AuthMiddleware,
	gameSessionService service.GameSessionService,
	roomService service.RoomService,
	chatService service.ChatService,
	participantRepo service.ParticipantLookup,
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
	client := websocket.NewClient(h.hub, conn, clientID)

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

// HandleMessage implements websocket.MessageHandler
func (h *WebSocketHandler) HandleMessage(client *websocket.Client, msg *websocket.Message) {
	switch msg.Type {
	case websocket.MessageTypeJoin:
		h.handleJoin(client, msg)
	case websocket.MessageTypeLeave:
		h.handleLeave(client)
	case websocket.MessageTypeChat:
		h.handleChat(client, msg)
	case websocket.MessageTypeAnswer:
		h.handleAnswer(client, msg)
	case websocket.MessageTypePing:
		h.handlePing(client)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (h *WebSocketHandler) handleJoin(client *websocket.Client, msg *websocket.Message) {
	var payload websocket.JoinPayload
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
	connectMsg, _ := websocket.NewMessage(websocket.MessageTypeConnect, map[string]interface{}{
		"clientId":      client.ID,
		"roomId":        payload.RoomID,
		"participantId": participantID,
		"isHost":        isHost,
	})
	client.SendMessage(connectMsg)

	// If game is active, send current state
	if h.gameSessionService.IsGameActive(payload.RoomID) {
		state := h.gameSessionService.GetGameState(payload.RoomID)
		stateMsg, _ := websocket.NewMessage(websocket.MessageTypeConnect, state)
		client.SendMessage(stateMsg)

		// Reconnect player in game session
		if participantID > 0 {
			h.gameSessionService.ReconnectPlayer(payload.RoomID, participantID, client.ID)
		}
	}

	log.Printf("Client %s joined room %d", client.ID, payload.RoomID)
}

func (h *WebSocketHandler) handleLeave(client *websocket.Client) {
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

func (h *WebSocketHandler) handleChat(client *websocket.Client, msg *websocket.Message) {
	var payload websocket.ChatPayload
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
	chatMsg, _ := websocket.NewMessage(websocket.MessageTypeChat, map[string]interface{}{
		"senderId":   client.UserID,
		"senderName": senderName,
		"content":    payload.Content,
	})
	h.hub.BroadcastToRoom(roomID, chatMsg)
}

func (h *WebSocketHandler) handleAnswer(client *websocket.Client, msg *websocket.Message) {
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
	ackMsg, _ := websocket.NewMessage(websocket.MessageTypeAnswerResult, map[string]interface{}{
		"questionId": payload.QuestionID,
		"received":   true,
	})
	client.SendMessage(ackMsg)
}

func (h *WebSocketHandler) handlePing(client *websocket.Client) {
	pongMsg, _ := websocket.NewMessage(websocket.MessageTypePong, nil)
	client.SendMessage(pongMsg)
}

func (h *WebSocketHandler) sendError(client *websocket.Client, code, message string) {
	errMsg, _ := websocket.NewMessage(websocket.MessageTypeError, websocket.ErrorPayload{
		Code:    code,
		Message: message,
	})
	client.SendMessage(errMsg)
}
