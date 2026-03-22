package response

// PageResponse represents a paginated response
type PageResponse[T any] struct {
	Content       []T   `json:"content"`
	Page          int   `json:"page"`
	Size          int   `json:"size"`
	TotalElements int64 `json:"totalElements"`
	TotalPages    int64 `json:"totalPages"`
}

// NewPageResponse creates a new page response
func NewPageResponse[T any](content []T, page, size int, totalElements int64) *PageResponse[T] {
	totalPages := (totalElements + int64(size) - 1) / int64(size)
	return &PageResponse[T]{
		Content:       content,
		Page:          page,
		Size:          size,
		TotalElements: totalElements,
		TotalPages:    totalPages,
	}
}
