package util

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Internal server error"`
	Message string `json:"message" example:"Failed to fetch data"`
}
