package models

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest represents the request body for creating a new user
// @Description Request model for user creation
type CreateUserRequest struct {
	// @Description User's full name
	Name string `json:"name" binding:"required" example:"John Smith"`
	
	// @Description User's email address (must be unique)
	Email string `json:"email" binding:"required,email" example:"john.smith@example.com"`
	
	// @Description User's password (minimum 6 characters)
	// @Required
	Password string `json:"password" binding:"required,min=6" example:"secretpassword123" minLength:"6"`
}

// UserResponse represents the response after user creation
// @Description Response model for user data
type UserResponse struct {
	// @Description Unique identifier for the user
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	
	// @Description User's full name
	Name string `json:"name" example:"John Smith"`
	
	// @Description User's email address
	Email string `json:"email" example:"john.smith@example.com"`
	
	// @Description When the user was created
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	
	// @Description When the user was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// User represents a user in the system
// @Description User model
type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string     `gorm:"size:100;not null" json:"name" example:"John Smith"`
	Email        string     `gorm:"size:255;not null;unique" json:"email" example:"john.smith@example.com"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" example:"2024-01-01T00:00:00Z"`
	Birthdays    []Birthday `gorm:"foreignKey:UserID" json:"-"` // Using json:"-" to exclude from Swagger docs
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// LoginRequest represents the request body for user login
// @Description Request model for user login
type LoginRequest struct {
	// @Description User's email address
	Email string `json:"email" binding:"required,email" example:"john.smith@example.com"`
	
	// @Description User's password
	Password string `json:"password" binding:"required" example:"secretpassword123"`
}

// LoginResponse represents the response for a successful login
// @Description Response model for successful login
type LoginResponse struct {
	// @Description JWT access token
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	
	// @Description Basic user information
	User UserResponse `json:"user"`
}

// UpdateUserRequest represents the request body for updating a user
// @Description Request model for updating user information
type UpdateUserRequest struct {
	// @Description User's full name
	Name string `json:"name" binding:"required" example:"John Smith"`
	
	// @Description User's email address (must be unique)
	Email string `json:"email" binding:"required,email" example:"john.smith@example.com"`
	
	// @Description User's new password (optional, minimum 6 characters if provided)
	Password string `json:"password" binding:"omitempty,min=6" example:"newpassword123" minLength:"6" swaggertype:"string"`
} 