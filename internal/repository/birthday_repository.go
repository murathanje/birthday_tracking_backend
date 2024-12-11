package repository

import (
	"github.com/google/uuid"
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"gorm.io/gorm"
)

type BirthdayRepository struct {
	db *gorm.DB
}

func NewBirthdayRepository(db *gorm.DB) *BirthdayRepository {
	return &BirthdayRepository{db: db}
}

func (r *BirthdayRepository) Create(birthday *models.Birthday) error {
	return r.db.Create(birthday).Error
}

func (r *BirthdayRepository) GetAll() ([]models.Birthday, error) {
	var birthdays []models.Birthday
	err := r.db.Preload("Category").Find(&birthdays).Error
	return birthdays, err
}

func (r *BirthdayRepository) GetByID(id uuid.UUID) (*models.Birthday, error) {
	var birthday models.Birthday
	err := r.db.Preload("Category").First(&birthday, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &birthday, nil
}

func (r *BirthdayRepository) GetByUserID(userID uuid.UUID) ([]models.Birthday, error) {
	var birthdays []models.Birthday
	err := r.db.Preload("Category").
		Where("user_id = ?", userID).
		Find(&birthdays).Error
	return birthdays, err
}

func (r *BirthdayRepository) Update(birthday *models.Birthday) error {
	return r.db.Save(birthday).Error
}

func (r *BirthdayRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Birthday{}, "id = ?", id).Error
}

func (r *BirthdayRepository) GetByCategory(categoryID uuid.UUID) ([]models.Birthday, error) {
	var birthdays []models.Birthday
	err := r.db.Preload("Category").
		Where("category_id = ?", categoryID).
		Find(&birthdays).Error
	return birthdays, err
} 