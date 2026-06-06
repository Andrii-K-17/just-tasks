package services

import (
	"errors"

	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/Andrii-K-17/just-tasks/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// ErrUsernameTaken is returned when a registration username is already in use.
var ErrUsernameTaken = errors.New("this username is already taken")

// ErrInvalidCredentials is returned on a failed login attempt.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthService handles user registration, login, and account management.
type AuthService struct {
	repo repository.UserRepository
}

// NewAuthService initializes and returns a new AuthService.
func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Register validates uniqueness, hashes the password, and creates a new user.
func (s *AuthService) Register(username, password string) (*models.User, error) {
	exists, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(username, string(hash))
}

// Login verifies credentials and returns the authenticated user.
func (s *AuthService) Login(username, password string) (*models.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// GetByID fetches a user by their primary key.
func (s *AuthService) GetByID(id int) (*models.User, error) {
	return s.repo.FindByID(id)
}

// DeleteAccount removes a user record from the database.
func (s *AuthService) DeleteAccount(userID int) error {
	return s.repo.Delete(userID)
}
