package service

import (
	"errors"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/dto/request"
	"github.com/huy/quizme-backend/internal/dto/response"
	"github.com/huy/quizme-backend/internal/repository"
	"gorm.io/gorm"
)

// ChatService handles chat-related operations
type ChatService interface {
	GetChatHistory(roomID uint, limit int) ([]*response.ChatMessageResponse, error)
	SendMessage(req *request.ChatMessageRequest, userID *uint) (*response.ChatMessageResponse, error)
}

type chatService struct {
	chatRepo repository.RoomChatRepository
	roomRepo repository.RoomRepository
}

// NewChatService creates a new chat service
func NewChatService(chatRepo repository.RoomChatRepository, roomRepo repository.RoomRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
		roomRepo: roomRepo,
	}
}

func (s *chatService) GetChatHistory(roomID uint, limit int) ([]*response.ChatMessageResponse, error) {
	// Verify room exists
	_, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}

	chats, err := s.chatRepo.FindByRoomID(roomID, limit)
	if err != nil {
		return nil, err
	}

	return response.FromRoomChats(chats), nil
}

func (s *chatService) SendMessage(req *request.ChatMessageRequest, userID *uint) (*response.ChatMessageResponse, error) {
	// Verify room exists
	_, err := s.roomRepo.FindByID(req.RoomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}

	isGuest := userID == nil
	chat := &domain.RoomChat{
		RoomID:    req.RoomID,
		UserID:    userID,
		IsGuest:   isGuest,
		GuestName: req.GuestName,
		Message:   req.Content,
	}

	if err := s.chatRepo.Create(chat); err != nil {
		return nil, err
	}

	return response.FromRoomChat(chat), nil
}
