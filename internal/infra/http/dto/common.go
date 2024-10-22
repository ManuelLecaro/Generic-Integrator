package dto

// ErrorResponseDTO represents the error response structure.
type ErrorResponseDTO struct {
	Message string `json:"message"` // Error message to describe the issue
}

// SuccessResponseDTO represents the success response structure.
type SuccessResponseDTO struct {
	Message string `json:"message"` // Success message indicating the result
}
