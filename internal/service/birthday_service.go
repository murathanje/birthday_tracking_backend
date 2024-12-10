package service

import (
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"github.com/murathanje/birthday_tracking_backend/internal/repository"
)

type BirthdayService struct {
	repo *repository.BirthdayRepository
}

func NewBirthdayService(repo *repository.BirthdayRepository) *BirthdayService {
	return &BirthdayService{repo: repo}
}

func (s *BirthdayService) CreateBirthday(birthday *models.Birthday) error {
	return s.repo.Create(birthday)
}

func (s *BirthdayService) GetAllBirthdays() ([]models.Birthday, error) {
	return s.repo.GetAll()
}

func (s *BirthdayService) GetBirthdayByID(id uint) (*models.Birthday, error) {
	return s.repo.GetByID(id)
}

func (s *BirthdayService) UpdateBirthday(birthday *models.Birthday) error {
	return s.repo.Update(birthday)
}

func (s *BirthdayService) DeleteBirthday(id uint) error {
	return s.repo.Delete(id)
} 