package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateBirthdayRequest represents the request for creating a birthday
// @Description Request model for creating a birthday record
type CreateBirthdayRequest struct {
	// @Description Name of the person
	Name string `json:"name" binding:"required" example:"John Doe"`

	// @Description Birthday date (format: MM-DD)
	BirthDate string `json:"birth_date" binding:"required" example:"05-15"`

	// @Description Category ID for grouping birthdays
	CategoryID uuid.UUID `json:"category_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440002"`

	// @Description Optional notes about the birthday
	Notes string `json:"notes,omitempty" example:"Best friend from college"`
}

// Birthday represents a birthday record
// @Description Birthday model for tracking birthdays
type Birthday struct {
	ID         uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID     uuid.UUID   `gorm:"type:uuid;not null" json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	User       User        `gorm:"foreignKey:UserID" json:"-"`
	Name       string      `gorm:"size:100;not null" json:"name" example:"John Doe"`
	BirthMonth int         `gorm:"not null" json:"birth_month" example:"5"`
	BirthDay   int         `gorm:"not null" json:"birth_day" example:"15"`
	CategoryID uuid.UUID   `gorm:"type:uuid;not null" json:"category_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Category   Category    `gorm:"foreignKey:CategoryID" json:"-"`
	Notes      string      `gorm:"type:text" json:"notes" example:"Best friend from college"`
	CreatedAt  time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// BirthdayResponse represents the response for birthday operations
// @Description Response model for birthday operations
type BirthdayResponse struct {
	// @Description Unique identifier for the birthday record
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// @Description User ID who owns this birthday record
	UserID uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001"`

	// @Description Name of the person
	Name string `json:"name" example:"John Doe"`

	// @Description Birthday date (format: MM-DD)
	BirthDate string `json:"birth_date" example:"05-15"`

	// @Description Category ID for grouping birthdays
	CategoryID uuid.UUID `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440002"`

	// @Description Optional notes about the birthday
	Notes string `json:"notes,omitempty" example:"Best friend from college"`

	// @Description When the record was created
	CreatedAt time.Time `json:"created_at"`

	// @Description When the record was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts Birthday model to BirthdayResponse
func (b *Birthday) ToResponse() *BirthdayResponse {
	return &BirthdayResponse{
		ID:        b.ID,
		UserID:    b.UserID,
		Name:      b.Name,
		BirthDate: fmt.Sprintf("%02d-%02d", b.BirthMonth, b.BirthDay),
		CategoryID: b.CategoryID,
		Notes:     b.Notes,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
} 