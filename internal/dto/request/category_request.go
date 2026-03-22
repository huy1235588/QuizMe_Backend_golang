package request

// CategoryRequest represents the request body for category operations
type CategoryRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Description *string `json:"description"`
	IconURL     *string `json:"iconUrl"`
	IsActive    *bool   `json:"isActive"`
}
