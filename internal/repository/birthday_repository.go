package repository

import (
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
	err := r.db.Find(&birthdays).Error
	return birthdays, err
}

func (r *BirthdayRepository) GetByID(id uint) (*models.Birthday, error) {
	var birthday models.Birthday
	err := r.db.First(&birthday, id).Error
	if err != nil {
		return nil, err
	}
	return &birthday, nil
}

func (r *BirthdayRepository) Update(birthday *models.Birthday) error {
	return r.db.Save(birthday).Error
}

func (r *BirthdayRepository) Delete(id uint) error {
	return r.db.Delete(&models.Birthday{}, id).Error
} 