package repository

import (
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"gorm.io/gorm"
)

type BirthdayRepository interface {
	Create(birthday *models.Birthday) error
	GetAll() ([]models.Birthday, error)
	GetByID(id uint) (*models.Birthday, error)
	Update(birthday *models.Birthday) error
	Delete(id uint) error
}

type birthdayRepository struct {
	db *gorm.DB
}

func NewBirthdayRepository(db *gorm.DB) BirthdayRepository {
	return &birthdayRepository{db: db}
}

func (r *birthdayRepository) Create(birthday *models.Birthday) error {
	return r.db.Create(birthday).Error
}

func (r *birthdayRepository) GetAll() ([]models.Birthday, error) {
	var birthdays []models.Birthday
	err := r.db.Find(&birthdays).Error
	return birthdays, err
}

func (r *birthdayRepository) GetByID(id uint) (*models.Birthday, error) {
	var birthday models.Birthday
	err := r.db.First(&birthday, id).Error
	return &birthday, err
}

func (r *birthdayRepository) Update(birthday *models.Birthday) error {
	return r.db.Save(birthday).Error
}

func (r *birthdayRepository) Delete(id uint) error {
	return r.db.Delete(&models.Birthday{}, id).Error
} 