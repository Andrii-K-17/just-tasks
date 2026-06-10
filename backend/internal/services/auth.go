package services

import (
	"errors"
	"time"

	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/Andrii-K-17/just-tasks/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// ErrUsernameTaken is returned when a registration username is already in use.
var ErrUsernameTaken = errors.New("this username is already taken")

// ErrInvalidCredentials is returned on a failed login attempt.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrInvalidRefreshToken is returned when a refresh token is missing, expired, or already used.
var ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")

// RefreshExpiry is the lifetime of a refresh token.
const RefreshExpiry = 30 * 24 * time.Hour

// TokenPair holds both the access JWT and the opaque refresh token.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// AuthService handles user registration, login, and account management.
type AuthService struct {
	repo        repository.UserRepository
	refreshRepo repository.RefreshTokenRepository
}

// NewAuthService initializes and returns a new AuthService.
func NewAuthService(repo repository.UserRepository, refreshRepo repository.RefreshTokenRepository) *AuthService {
	return &AuthService{repo: repo, refreshRepo: refreshRepo}
}

// issueTokenPair creates an access JWT and a refresh token for the given user.
func (s *AuthService) issueTokenPair(userID int, jwtSecret string, jwtExpiry time.Duration) (*TokenPair, error) {
	accessToken, err := IssueJWT(userID, jwtSecret, jwtExpiry)
	if err != nil {
		return nil, err
	}

	rawRefresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenHash := HashRefreshToken(rawRefresh)
	expiresAt := time.Now().Add(RefreshExpiry)

	if _, err := s.refreshRepo.Create(userID, tokenHash, expiresAt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}

// Register validates uniqueness, hashes the password, and creates a new user.
func (s *AuthService) Register(username, password, jwtSecret string, jwtExpiry time.Duration) (*models.User, *TokenPair, error) {
	exists, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.repo.Create(username, string(hash))
	if err != nil {
		return nil, nil, err
	}

	pair, err := s.issueTokenPair(user.ID, jwtSecret, jwtExpiry)
	if err != nil {
		return nil, nil, err
	}

	return user, pair, nil
}

// Login verifies credentials and returns the authenticated user with a token pair.
func (s *AuthService) Login(username, password, jwtSecret string, jwtExpiry time.Duration) (*models.User, *TokenPair, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password),
	); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	pair, err := s.issueTokenPair(user.ID, jwtSecret, jwtExpiry)
	if err != nil {
		return nil, nil, err
	}

	return user, pair, nil
}

// Refresh validates the incoming refresh token, rotates it, and returns a new token pair.
func (s *AuthService) Refresh(rawToken, jwtSecret string, jwtExpiry time.Duration) (*TokenPair, error) {
	tokenHash := HashRefreshToken(rawToken)

	stored, err := s.refreshRepo.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if err := s.refreshRepo.DeleteByTokenHash(tokenHash); err != nil {
		return nil, err
	}

	if time.Now().After(stored.ExpiresAt) {
		return nil, ErrInvalidRefreshToken
	}

	return s.issueTokenPair(stored.UserID, jwtSecret, jwtExpiry)
}

// Logout removes the given refresh token from the store.
func (s *AuthService) Logout(rawToken string) error {
	if rawToken == "" {
		return nil
	}
	tokenHash := HashRefreshToken(rawToken)

	return s.refreshRepo.DeleteByTokenHash(tokenHash)
}

// GetByID fetches a user by their primary key.
func (s *AuthService) GetByID(id int) (*models.User, error) {
	return s.repo.FindByID(id)
}

// DeleteAccount removes a user record and all associated refresh tokens from the database.
func (s *AuthService) DeleteAccount(userID int) error {
	if err := s.refreshRepo.DeleteAllByUserID(userID); err != nil {
		return err
	}
	return s.repo.Delete(userID)
}
