package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to rooms
type Hub struct {
	// Registered clients by room
	rooms map[uint]map[*Client]bool

	// Clients by ID for quick lookup
	clients map[string]*Client

	// Inbound messages from clients
	Broadcast chan *BroadcastMessage

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Message handler
	messageHandler MessageHandler

	mu sync.RWMutex
}

// BroadcastMessage is a message to broadcast to a room
type BroadcastMessage struct {
	RoomID  uint
	Message *Message
	Exclude *Client // Client to exclude from broadcast
}

// MessageHandler handles incoming WebSocket messages
type MessageHandler interface {
	HandleMessage(client *Client, msg *Message)
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uint]map[*Client]bool),
		clients:    make(map[string]*Client),
		Broadcast:  make(chan *BroadcastMessage, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// SetMessageHandler sets the message handler for the hub
func (h *Hub) SetMessageHandler(handler MessageHandler) {
	h.messageHandler = handler
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case broadcast := <-h.Broadcast:
			h.broadcastToRoom(broadcast)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client.ID] = client
	log.Printf("Client %s registered", client.ID)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients, client.ID)

		// Remove from room if in one
		if client.RoomID > 0 {
			if room, ok := h.rooms[client.RoomID]; ok {
				delete(room, client)
				if len(room) == 0 {
					delete(h.rooms, client.RoomID)
				}
			}
		}

		log.Printf("Client %s unregistered from room %d", client.ID, client.RoomID)
	}
}

func (h *Hub) broadcastToRoom(broadcast *BroadcastMessage) {
	h.mu.RLock()
	room, ok := h.rooms[broadcast.RoomID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	data, err := json.Marshal(broadcast.Message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	h.mu.RLock()
	for client := range room {
		if broadcast.Exclude != nil && client == broadcast.Exclude {
			continue
		}
		select {
		case client.Send <- data:
		default:
			// Channel full, will be cleaned up by WritePump
		}
	}
	h.mu.RUnlock()
}

// JoinRoom adds a client to a room
func (h *Hub) JoinRoom(client *Client, roomID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Remove from previous room if any
	if client.RoomID > 0 && client.RoomID != roomID {
		if room, ok := h.rooms[client.RoomID]; ok {
			delete(room, client)
			if len(room) == 0 {
				delete(h.rooms, client.RoomID)
			}
		}
	}

	// Add to new room
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
	client.RoomID = roomID

	log.Printf("Client %s joined room %d", client.ID, roomID)
}

// LeaveRoom removes a client from their current room
func (h *Hub) LeaveRoom(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client.RoomID > 0 {
		if room, ok := h.rooms[client.RoomID]; ok {
			delete(room, client)
			if len(room) == 0 {
				delete(h.rooms, client.RoomID)
			}
		}
		log.Printf("Client %s left room %d", client.ID, client.RoomID)
		client.RoomID = 0
	}
}

// BroadcastToRoom sends a message to all clients in a room
func (h *Hub) BroadcastToRoom(roomID uint, msg *Message) {
	h.Broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: msg,
	}
}

// BroadcastToRoomExcept sends a message to all clients in a room except one
func (h *Hub) BroadcastToRoomExcept(roomID uint, msg *Message, exclude *Client) {
	h.Broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: msg,
		Exclude: exclude,
	}
}

// SendToClient sends a message to a specific client
func (h *Hub) SendToClient(clientID string, msg *Message) {
	h.mu.RLock()
	client, ok := h.clients[clientID]
	h.mu.RUnlock()

	if ok {
		client.SendMessage(msg)
	}
}

// SendToUser sends a message to all connections of a user
func (h *Hub) SendToUser(userID uint, msg *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.UserID != nil && *client.UserID == userID {
			client.SendMessage(msg)
		}
	}
}

// GetRoomClients returns all clients in a room
func (h *Hub) GetRoomClients(roomID uint) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var clients []*Client
	if room, ok := h.rooms[roomID]; ok {
		for client := range room {
			clients = append(clients, client)
		}
	}
	return clients
}

// GetRoomClientCount returns the number of clients in a room
func (h *Hub) GetRoomClientCount(roomID uint) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.rooms[roomID]; ok {
		return len(room)
	}
	return 0
}

// GetClient returns a client by ID
func (h *Hub) GetClient(clientID string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[clientID]
}

// HandleMessage routes messages to the appropriate handler
func (h *Hub) HandleMessage(client *Client, msg *Message) {
	if h.messageHandler != nil {
		h.messageHandler.HandleMessage(client, msg)
	}
}
