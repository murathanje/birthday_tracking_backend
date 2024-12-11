package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/murathanje/birthday_tracking_backend/internal/config"
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"github.com/murathanje/birthday_tracking_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	tokenExpiration      = 24 * time.Hour
)

type UserService struct {
	repo   *repository.UserRepository
	config *config.Config
}

func NewUserService(repo *repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		repo:   repo,
		config: cfg,
	}
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user with hashed password
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UserService) UpdateUser(id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and if it's already taken
	if user.Email != req.Email {
		existingUser, err := s.repo.GetByEmail(req.Email)
		if err == nil && existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	// Update basic information
	user.Name = req.Name
	user.Email = req.Email

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Save changes
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(tokenExpiration).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Create response
	return &models.LoginResponse{
		Token: tokenString,
		User:  *user.ToResponse(),
	}, nil
}

// GetJWTSecret returns the JWT signing secret
func (s *UserService) GetJWTSecret() []byte {
	return []byte(s.config.JWTSecret)
} 