package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"github.com/murathanje/birthday_tracking_backend/internal/repository"
)

type BirthdayService struct {
	repo *repository.BirthdayRepository
}

func NewBirthdayService(repo *repository.BirthdayRepository) *BirthdayService {
	return &BirthdayService{repo: repo}
}

func (s *BirthdayService) CreateBirthday(userID uuid.UUID, req *models.CreateBirthdayRequest) (*models.Birthday, error) {
	// Parse birth date (MM-DD format)
	parts := strings.Split(req.BirthDate, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid birth date format, expected MM-DD")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month")
	}

	day, err := strconv.Atoi(parts[1])
	if err != nil || day < 1 || day > 31 {
		return nil, fmt.Errorf("invalid day")
	}

	daysInMonth := getDaysInMonth(month)
	if day > daysInMonth {
		return nil, fmt.Errorf("invalid day for month %d", month)
	}

	birthday := &models.Birthday{
		UserID:     userID,
		Name:       req.Name,
		BirthMonth: month,
		BirthDay:   day,
		Category:   req.Category,
		Notes:      req.Notes,
	}

	if err := s.repo.Create(birthday); err != nil {
		return nil, err
	}

	return birthday, nil
}

func (s *BirthdayService) GetByID(id uuid.UUID) (*models.Birthday, error) {
	return s.repo.GetByID(id)
}

func (s *BirthdayService) GetByUserID(userID uuid.UUID) ([]models.Birthday, error) {
	return s.repo.GetByUserID(userID)
}

func (s *BirthdayService) Update(birthday *models.Birthday) error {
	return s.repo.Update(birthday)
}

func (s *BirthdayService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *BirthdayService) GetByCategory(category string) ([]models.Birthday, error) {
	return s.repo.GetByCategory(category)
}

// Helper function to get days in a month
func getDaysInMonth(month int) int {
	switch month {
	case 4, 6, 9, 11:
		return 30
	case 2:
		return 29
	default:
		return 31
	}
} 