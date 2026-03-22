package response

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse[T any] struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// PageResponse represents a paginated response
type PageResponse[T any] struct {
	Content       T     `json:"content"`
	Page          int   `json:"page"`
	Size          int   `json:"size"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int64 `json:"totalPages"`
}

// Error returns an error response
func Error(message string) ErrorResponse {
	return ErrorResponse{
		Status:  false,
		Message: message,
	}
}

// Success returns a success response
func Success[T any](data T, message string) SuccessResponse[T] {
	return SuccessResponse[T]{
		Status:  true,
		Message: message,
		Data:    data,
	}
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID        uint    `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
	IsActive  bool    `json:"isActive"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// ChatMessageResponse represents a chat message in responses
type ChatMessageResponse struct {
	ID        uint    `json:"id"`
	RoomID    uint    `json:"roomId"`
	UserID    *uint   `json:"userId,omitempty"`
	Username  string  `json:"username"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"createdAt"`
}

// GameResultResponse represents game result data in responses
type GameResultResponse struct {
	ID             uint    `json:"id"`
	RoomID         uint    `json:"roomId"`
	UserID         uint    `json:"userId"`
	Score          int     `json:"score"`
	CorrectAnswers int     `json:"correctAnswers"`
	CreatedAt      string  `json:"createdAt"`
}

// FromUser converts a user domain model to UserResponse
func FromUser(user interface{}) *UserResponse {
	// This is a placeholder that will need to be implemented
	// when we have the actual user domain model
	return nil
}
