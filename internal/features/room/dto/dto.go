package dto

import (
	"time"

	"github.com/huy/quizme-backend/internal/domain/enums"
	"github.com/huy/quizme-backend/internal/dto/response"
	quizdto "github.com/huy/quizme-backend/internal/features/quiz/dto"
	"github.com/huy/quizme-backend/internal/features/room/domain"
	userdto "github.com/huy/quizme-backend/internal/features/user/dto"
)

// ==== REQUEST DTOs ====

// RoomRequest represents the request body for room creation
type RoomRequest struct {
	Name       string  `json:"name" validate:"required,max=100"`
	QuizID     uint    `json:"quizId" validate:"required"`
	MaxPlayers int     `json:"maxPlayers" validate:"min=2,max=100"`
	Password   *string `json:"password"`
	IsPublic   *bool   `json:"isPublic"`
}

// JoinRoomRequest represents the request body for joining a room by code
type JoinRoomRequest struct {
	Code      string  `json:"code" validate:"required"`
	GuestName *string `json:"guestName"`
	Password  *string `json:"password"`
}

// JoinRoomByIDRequest represents the request body for joining a room by ID
type JoinRoomByIDRequest struct {
	GuestName *string `json:"guestName"`
	Password  *string `json:"password"`
}

// ChatMessageRequest represents the request body for sending a chat message
type ChatMessageRequest struct {
	RoomID    uint    `json:"roomId" validate:"required"`
	Content   string  `json:"content" validate:"required,max=500"`
	GuestName *string `json:"guestName"`
}

// UpdateRoomRequest represents the request body for updating a room
type UpdateRoomRequest struct {
	Name       *string `json:"name" validate:"omitempty,max=100"`
	MaxPlayers *int    `json:"maxPlayers" validate:"omitempty,min=2,max=100"`
	Password   *string `json:"password"`
	IsPublic   *bool   `json:"isPublic"`
}

// ==== RESPONSE DTOs ====

// RoomResponse represents a room in API responses
type RoomResponse struct {
	ID                 uint                   `json:"id"`
	Name               string                 `json:"name"`
	Code               string                 `json:"code"`
	QuizID             uint                   `json:"quizId"`
	HostID             uint                   `json:"hostId"`
	Quiz               *quizdto.QuizResponse  `json:"quiz,omitempty"`
	Host               *response.UserResponse `json:"host,omitempty"`
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
	ID          uint                   `json:"id"`
	User        *response.UserResponse `json:"user,omitempty"`
	Score       int                    `json:"score"`
	IsHost      bool                   `json:"isHost"`
	JoinedAt    string                 `json:"joinedAt"`
	LeftAt      *string                `json:"leftAt,omitempty"`
	IsGuest     bool                   `json:"isGuest"`
	GuestName   *string                `json:"guestName,omitempty"`
	DisplayName string                 `json:"displayName"`
}

// ChatMessageResponse represents a chat message in API responses
type ChatMessageResponse struct {
	ID          uint                   `json:"id"`
	RoomID      uint                   `json:"roomId"`
	User        *response.UserResponse `json:"user,omitempty"`
	IsGuest     bool                   `json:"isGuest"`
	GuestName   *string                `json:"guestName,omitempty"`
	Message     string                 `json:"message"`
	SentAt      string                 `json:"sentAt"`
	DisplayName string                 `json:"displayName"`
}

// ==== CONVERSION FUNCTIONS ====

// userDTOToResponse converts a user DTO UserResponse to response UserResponse
func userDTOToResponse(userResp *userdto.UserResponse) *response.UserResponse {
	if userResp == nil {
		return nil
	}
	return &response.UserResponse{
		ID:        userResp.ID,
		Username:  userResp.Username,
		Email:     userResp.Email,
		FullName:  userResp.FullName,
		AvatarURL: userResp.ProfileImage,
		IsActive:  userResp.IsActive,
		CreatedAt: userResp.CreatedAt,
		UpdatedAt: userResp.UpdatedAt,
	}
}

// FromRoom converts a Room domain model to RoomResponse
func FromRoom(room *domain.Room) *RoomResponse {
	resp := &RoomResponse{
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
			resp.CurrentPlayerCount++
		}
	}

	// Add quiz info
	if room.Quiz != nil {
		resp.Quiz = quizdto.FromQuiz(room.Quiz)
	}

	// Add host info
	if room.Host != nil {
		resp.Host = userDTOToResponse(userdto.FromUser(room.Host))
	}

	// Add time info
	if room.StartTime != nil {
		t := room.StartTime.Format(time.RFC3339)
		resp.StartTime = &t
	}
	if room.EndTime != nil {
		t := room.EndTime.Format(time.RFC3339)
		resp.EndTime = &t
	}

	// Add participants
	if len(room.Participants) > 0 {
		resp.Participants = make([]*ParticipantResponse, len(room.Participants))
		for i, p := range room.Participants {
			resp.Participants[i] = FromRoomParticipant(&p)
		}
	}

	return resp
}

// FromRoomParticipant converts a RoomParticipant domain model to ParticipantResponse
func FromRoomParticipant(participant *domain.RoomParticipant) *ParticipantResponse {
	resp := &ParticipantResponse{
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
		resp.LeftAt = &t
	}

	if participant.User != nil {
		resp.User = userDTOToResponse(userdto.FromUser(participant.User))
	}

	return resp
}

// FromRoomChat converts a RoomChat domain model to ChatMessageResponse
func FromRoomChat(chat *domain.RoomChat) *ChatMessageResponse {
	resp := &ChatMessageResponse{
		ID:        chat.ID,
		RoomID:    chat.RoomID,
		IsGuest:   chat.IsGuest,
		GuestName: chat.GuestName,
		Message:   chat.Message,
		SentAt:    chat.SentAt.Format(time.RFC3339),
	}

	if chat.User != nil {
		resp.User = userDTOToResponse(userdto.FromUser(chat.User))
		resp.DisplayName = chat.User.Username
	} else if chat.GuestName != nil {
		resp.DisplayName = *chat.GuestName
	} else {
		resp.DisplayName = "Guest"
	}

	return resp
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
