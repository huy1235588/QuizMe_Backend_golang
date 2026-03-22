package response

// ApiResponse represents the standard API response structure
type ApiResponse[T any] struct {
	Status  string `json:"status"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message"`
}

// Success creates a success response with data
func Success[T any](data T, message string) ApiResponse[T] {
	return ApiResponse[T]{
		Status:  "success",
		Data:    data,
		Message: message,
	}
}

// Created creates a success response for resource creation
func Created[T any](data T, message string) ApiResponse[T] {
	return ApiResponse[T]{
		Status:  "success",
		Data:    data,
		Message: message,
	}
}

// Error creates an error response
func Error(message string) ApiResponse[any] {
	return ApiResponse[any]{
		Status:  "error",
		Data:    nil,
		Message: message,
	}
}

// ErrorWithData creates an error response with additional data
func ErrorWithData[T any](data T, message string) ApiResponse[T] {
	return ApiResponse[T]{
		Status:  "error",
		Data:    data,
		Message: message,
	}
}
