package request

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
