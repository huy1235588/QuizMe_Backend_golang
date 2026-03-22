package service

import (
	roomdto "github.com/huy/quizme-backend/internal/features/room/dto"
)

// RoomService handles room-related operations
type RoomService interface {
	CreateRoom(hostID uint, req *roomdto.RoomRequest) (*roomdto.RoomResponse, error)
	GetRoomByCode(code string) (*roomdto.RoomResponse, error)
	GetRoomByID(id uint) (*roomdto.RoomResponse, error)
	GetWaitingRooms() ([]*roomdto.RoomResponse, error)
	GetAvailableRooms(search *string, quizID *uint, page, pageSize int) ([]*roomdto.RoomResponse, int64, error)
	JoinRoom(roomID uint, userID *uint, guestName, password *string) (*roomdto.ParticipantResponse, error)
	JoinRoomByCode(code string, userID *uint, guestName, password *string) (*roomdto.RoomResponse, *roomdto.ParticipantResponse, error)
	LeaveRoom(roomID uint, userID *uint, guestName *string) error
	CloseRoom(roomID, hostID uint) error
	UpdateRoom(roomID, hostID uint, req *roomdto.UpdateRoomRequest) (*roomdto.RoomResponse, error)
	StartGame(roomID, hostID uint) (*roomdto.RoomResponse, error)
}

// ChatService handles chat-related operations
type ChatService interface {
	GetChatHistory(roomID uint, limit int) ([]*roomdto.ChatMessageResponse, error)
	SendMessage(req *roomdto.ChatMessageRequest, userID *uint) (*roomdto.ChatMessageResponse, error)
}
