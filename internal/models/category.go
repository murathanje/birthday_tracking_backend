package models

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a birthday category
// @Description Category model for grouping birthdays
type Category struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string     `gorm:"size:50;not null;unique" json:"name" example:"Family"`
	Icon      string     `gorm:"size:10;not null" json:"icon" example:"👨‍👩‍👧‍👦"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Birthdays []Birthday `gorm:"foreignKey:CategoryID" json:"-"` // Using json:"-" to exclude from Swagger docs
} 