package repository

import (
	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/features/room/domain"
)

// RoomRepository defines the interface for room data access
type RoomRepository interface {
	Create(room *domain.Room) error
	FindByID(id uint) (*domain.Room, error)
	FindByCode(code string) (*domain.Room, error)
	FindByHostID(hostID uint) ([]domain.Room, error)
	FindWaiting() ([]domain.Room, error)
	FindAvailable(search *string, quizID *uint, page, pageSize int) ([]domain.Room, int64, error)
	Update(room *domain.Room) error
	Delete(id uint) error
	UpdateStatus(id uint, status enums.RoomStatus) error
}

// RoomParticipantRepository defines the interface for room participant data access
type RoomParticipantRepository interface {
	Create(participant *domain.RoomParticipant) error
	FindByID(id uint) (*domain.RoomParticipant, error)
	FindByRoomID(roomID uint) ([]domain.RoomParticipant, error)
	FindByRoomAndUser(roomID, userID uint) (*domain.RoomParticipant, error)
	Update(participant *domain.RoomParticipant) error
	Delete(id uint) error
	DeleteByRoomID(roomID uint) error
	CountByRoomID(roomID uint) (int64, error)
	UpdateScore(id uint, score int) error
}

// RoomChatRepository defines the interface for room chat data access
type RoomChatRepository interface {
	Create(chat *domain.RoomChat) error
	FindByRoomID(roomID uint, limit int) ([]domain.RoomChat, error)
	DeleteByRoomID(roomID uint) error
}
