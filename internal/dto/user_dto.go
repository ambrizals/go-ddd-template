package dto

// RegisterRequest is the input for user registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"strongpassword123"`
	FullName string `json:"full_name" validate:"required" example:"John Doe"`
}

// UserResponse is the output for user data.
type UserResponse struct {
	ID        uint   `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	FullName  string `json:"full_name" example:"John Doe"`
	CreatedAt string `json:"created_at" example:"2023-10-27T10:00:00Z"`
}

// ErrorResponse represents a standardized error message.
type ErrorResponse struct {
	Error   string `json:"error" example:"Validation failed"`
	Message string `json:"message" example:"Email is required"`
}
