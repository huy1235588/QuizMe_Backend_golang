package response

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain"
	"github.com/huy/quizme-backend/internal/domain/enums"
)

// RoomResponse represents a room in API responses
type RoomResponse struct {
	ID                 uint                   `json:"id"`
	Name               string                 `json:"name"`
	Code               string                 `json:"code"`
	QuizID             uint                   `json:"quizId"`
	HostID             uint                   `json:"hostId"`
	Quiz               *QuizResponse          `json:"quiz,omitempty"`
	Host               *UserResponse          `json:"host,omitempty"`
	HasPassword        bool                   `json:"hasPassword"`
	IsPublic           bool                   `json:"isPublic"`
	CurrentPlayerCount int                    `json:"currentPlayerCount"`
	MaxPlayers         int                    `json:"maxPlayers"`
	Status             enums.RoomStatus       `json:"status"`
	StartTime          *string                `json:"startTime,omitempty"`
	EndTime            *string                `json:"endTime,omitempty"`
	CreatedAt          string                 `json:"createdAt"`
	Participants       []*ParticipantResponse `json:"participants,omitempty"`
}

// ParticipantResponse represents a room participant in API responses
type ParticipantResponse struct {
	ID          uint          `json:"id"`
	User        *UserResponse `json:"user,omitempty"`
	Score       int           `json:"score"`
	IsHost      bool          `json:"isHost"`
	JoinedAt    string        `json:"joinedAt"`
	LeftAt      *string       `json:"leftAt,omitempty"`
	IsGuest     bool          `json:"isGuest"`
	GuestName   *string       `json:"guestName,omitempty"`
	DisplayName string        `json:"displayName"`
}

// ChatMessageResponse represents a chat message in API responses
type ChatMessageResponse struct {
	ID          uint          `json:"id"`
	RoomID      uint          `json:"roomId"`
	User        *UserResponse `json:"user,omitempty"`
	IsGuest     bool          `json:"isGuest"`
	GuestName   *string       `json:"guestName,omitempty"`
	Message     string        `json:"message"`
	SentAt      string        `json:"sentAt"`
	DisplayName string        `json:"displayName"`
}

// FromRoom converts a Room domain model to RoomResponse
func FromRoom(room *domain.Room) *RoomResponse {
	response := &RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Code:        room.Code,
		QuizID:      room.QuizID,
		HostID:      room.HostID,
		HasPassword: room.HasPassword(),
		IsPublic:    room.IsPublic,
		MaxPlayers:  room.MaxPlayers,
		Status:      room.Status,
		CreatedAt:   room.CreatedAt.Format(time.RFC3339),
	}

	// Count current players
	for _, p := range room.Participants {
		if p.LeftAt == nil {
			response.CurrentPlayerCount++
		}
	}

	// Add quiz info
	if room.Quiz != nil {
		response.Quiz = FromQuiz(room.Quiz)
	}

	// Add host info
	if room.Host != nil {
		response.Host = FromUser(room.Host)
	}

	// Add time info
	if room.StartTime != nil {
		t := room.StartTime.Format(time.RFC3339)
		response.StartTime = &t
	}
	if room.EndTime != nil {
		t := room.EndTime.Format(time.RFC3339)
		response.EndTime = &t
	}

	// Add participants
	if len(room.Participants) > 0 {
		response.Participants = make([]*ParticipantResponse, len(room.Participants))
		for i, p := range room.Participants {
			response.Participants[i] = FromRoomParticipant(&p)
		}
	}

	return response
}

// FromRoomParticipant converts a RoomParticipant domain model to ParticipantResponse
func FromRoomParticipant(participant *domain.RoomParticipant) *ParticipantResponse {
	response := &ParticipantResponse{
		ID:          participant.ID,
		Score:       participant.Score,
		IsHost:      participant.IsHost,
		JoinedAt:    participant.JoinedAt.Format(time.RFC3339),
		IsGuest:     participant.IsGuest,
		GuestName:   participant.GuestName,
		DisplayName: participant.GetDisplayName(),
	}

	if participant.LeftAt != nil {
		t := participant.LeftAt.Format(time.RFC3339)
		response.LeftAt = &t
	}

	if participant.User != nil {
		response.User = FromUser(participant.User)
	}

	return response
}

// FromRoomChat converts a RoomChat domain model to ChatMessageResponse
func FromRoomChat(chat *domain.RoomChat) *ChatMessageResponse {
	response := &ChatMessageResponse{
		ID:        chat.ID,
		RoomID:    chat.RoomID,
		IsGuest:   chat.IsGuest,
		GuestName: chat.GuestName,
		Message:   chat.Message,
		SentAt:    chat.SentAt.Format(time.RFC3339),
	}

	if chat.User != nil {
		response.User = FromUser(chat.User)
		response.DisplayName = chat.User.Username
	} else if chat.GuestName != nil {
		response.DisplayName = *chat.GuestName
	} else {
		response.DisplayName = "Guest"
	}

	return response
}

// FromRooms converts a slice of Room domain models to RoomResponse slice
func FromRooms(rooms []domain.Room) []*RoomResponse {
	responses := make([]*RoomResponse, len(rooms))
	for i, room := range rooms {
		responses[i] = FromRoom(&room)
	}
	return responses
}

// FromRoomChats converts a slice of RoomChat domain models to ChatMessageResponse slice
func FromRoomChats(chats []domain.RoomChat) []*ChatMessageResponse {
	responses := make([]*ChatMessageResponse, len(chats))
	for i, chat := range chats {
		responses[i] = FromRoomChat(&chat)
	}
	return responses
}
