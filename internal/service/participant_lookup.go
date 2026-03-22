package service

import "github.com/huy/quizme-backend/internal/domain"

// ParticipantLookup provides participant lookup functionality
type ParticipantLookup interface {
	FindByRoomAndUser(roomID, userID uint) (*domain.RoomParticipant, error)
}
