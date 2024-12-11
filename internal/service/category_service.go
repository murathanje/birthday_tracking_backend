package service

import (
	"github.com/google/uuid"
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"github.com/murathanje/birthday_tracking_backend/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) GetCategoryByName(name string) (*models.Category, error) {
	return s.repo.GetByName(name)
} 