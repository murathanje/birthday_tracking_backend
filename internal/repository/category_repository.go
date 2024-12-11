package repository

import (
	"github.com/google/uuid"
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
} 