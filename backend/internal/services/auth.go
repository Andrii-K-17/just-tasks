package services

import (
	"errors"

	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// ErrUsernameTaken is returned when a registration username is already in use.
var ErrUsernameTaken = errors.New("this username is already taken")

// ErrInvalidCredentials is returned on a failed login attempt.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthService handles user registration, login, and account management.
type AuthService struct {
	db *sqlx.DB
}

// NewAuthService initializes and returns a new AuthService.
func NewAuthService(db *sqlx.DB) *AuthService {
	return &AuthService{db: db}
}

// Register validates uniqueness, hashes the password, and creates a new user.
func (s *AuthService) Register(username, password string) (*models.User, error) {
	var exists bool
	if err := s.db.Get(&exists,
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username,
	); err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.db.QueryRowx(
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)
		 RETURNING id, username, created_at`,
		username, string(hash),
	).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Login verifies credentials and returns the authenticated user.
func (s *AuthService) Login(username, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Get(&user,
		"SELECT id, username, password_hash FROM users WHERE username=$1", username,
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

// GetByID fetches a user by their primary key.
func (s *AuthService) GetByID(id int) (*models.User, error) {
	var user models.User
	if err := s.db.Get(&user,
		"SELECT id, username FROM users WHERE id=$1", id,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteAccount removes a user record from the database.
func (s *AuthService) DeleteAccount(userID int) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id=$1", userID)
	return err
}
