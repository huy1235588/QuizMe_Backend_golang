package service

import (
	"errors"

	roomdto "github.com/huy/quizme-backend/internal/features/room/dto"
	roomdomain "github.com/huy/quizme-backend/internal/features/room/domain"
	roomrepo "github.com/huy/quizme-backend/internal/features/room/repository"
	"gorm.io/gorm"
)

type chatService struct {
	chatRepo roomrepo.RoomChatRepository
	roomRepo roomrepo.RoomRepository
}

// NewChatService creates a new chat service
func NewChatService(chatRepo roomrepo.RoomChatRepository, roomRepo roomrepo.RoomRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
		roomRepo: roomRepo,
	}
}

func (s *chatService) GetChatHistory(roomID uint, limit int) ([]*roomdto.ChatMessageResponse, error) {
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

	return roomdto.FromRoomChats(chats), nil
}

func (s *chatService) SendMessage(req *roomdto.ChatMessageRequest, userID *uint) (*roomdto.ChatMessageResponse, error) {
	// Verify room exists
	_, err := s.roomRepo.FindByID(req.RoomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}

	isGuest := userID == nil
	chat := &roomdomain.RoomChat{
		RoomID:    req.RoomID,
		UserID:    userID,
		IsGuest:   isGuest,
		GuestName: req.GuestName,
		Message:   req.Content,
	}

	if err := s.chatRepo.Create(chat); err != nil {
		return nil, err
	}

	return roomdto.FromRoomChat(chat), nil
}
