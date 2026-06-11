package services_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Andrii-K-17/just-tasks/internal/mocks"
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/Andrii-K-17/just-tasks/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const testSecret = "test_secret"
const testExpiry = time.Hour
const testRefreshExpiry = 7 * 24 * time.Hour

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenRepository)
	svc := services.NewAuthService(userRepo, refreshRepo)

	expected := &models.User{ID: 1, Username: "user1", CreatedAt: time.Now()}
	stored := &models.RefreshToken{ID: 1, UserID: 1, ExpiresAt: time.Now().Add(testRefreshExpiry)}

	userRepo.On("ExistsByUsername", "user1").Return(false, nil)
	userRepo.On("Create", "user1", mock.MatchedBy(func(s string) bool {
		return strings.HasPrefix(s, "$2a$")
	})).Return(expected, nil)
	refreshRepo.On("Create", 1, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
		Return(stored, nil)

	user, pair, err := svc.Register("user1", "password123", testSecret, testExpiry, testRefreshExpiry)

	require.NoError(t, err)
	assert.Equal(t, expected.Username, user.Username)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	userRepo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}

func TestAuthService_Register_UsernameTaken(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenRepository)
	svc := services.NewAuthService(userRepo, refreshRepo)

	userRepo.On("ExistsByUsername", "user1").Return(true, nil)

	_, _, err := svc.Register("user1", "password123", testSecret, testExpiry, testRefreshExpiry)

	assert.ErrorIs(t, err, services.ErrUsernameTaken)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenRepository)
	svc := services.NewAuthService(userRepo, refreshRepo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	stored := &models.User{ID: 1, Username: "user1", PasswordHash: string(hash)}
	storedToken := &models.RefreshToken{ID: 1, UserID: 1, ExpiresAt: time.Now().Add(testRefreshExpiry)}

	userRepo.On("FindByUsername", "user1").Return(stored, nil)
	refreshRepo.On("Create", 1, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
		Return(storedToken, nil)

	user, pair, err := svc.Login("user1", "password123", testSecret, testExpiry, testRefreshExpiry)

	require.NoError(t, err)
	assert.Equal(t, stored.ID, user.ID)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	userRepo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenRepository)
	svc := services.NewAuthService(userRepo, refreshRepo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	stored := &models.User{ID: 1, Username: "user1", PasswordHash: string(hash)}

	userRepo.On("FindByUsername", "user1").Return(stored, nil)

	_, _, err := svc.Login("user1", "wrong", testSecret, testExpiry, testRefreshExpiry)

	assert.ErrorIs(t, err, services.ErrInvalidCredentials)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenRepository)
	svc := services.NewAuthService(userRepo, refreshRepo)

	userRepo.On("FindByUsername", "ghost").Return(nil, errors.New("not found"))

	_, _, err := svc.Login("ghost", "password123", testSecret, testExpiry, testRefreshExpiry)

	assert.ErrorIs(t, err, services.ErrInvalidCredentials)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Refresh_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenRepository)
	svc := services.NewAuthService(userRepo, refreshRepo)

	rawToken := "some_raw_token"
	tokenHash := services.HashRefreshToken(rawToken)
	storedToken := &models.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	newToken := &models.RefreshToken{ID: 2, UserID: 1, ExpiresAt: time.Now().Add(testRefreshExpiry)}

	refreshRepo.On("FindByTokenHash", tokenHash).Return(storedToken, nil)
	refreshRepo.On("DeleteByTokenHash", tokenHash).Return(nil)
	refreshRepo.On("Create", 1, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
		Return(newToken, nil)

	pair, err := svc.Refresh(rawToken, testSecret, testExpiry, testRefreshExpiry)

	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	refreshRepo.AssertExpectations(t)
}

func TestAuthService_Refresh_Expired(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenRepository)
	svc := services.NewAuthService(userRepo, refreshRepo)

	rawToken := "expired_token"
	tokenHash := services.HashRefreshToken(rawToken)
	storedToken := &models.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	refreshRepo.On("FindByTokenHash", tokenHash).Return(storedToken, nil)
	refreshRepo.On("DeleteByTokenHash", tokenHash).Return(nil)

	_, err := svc.Refresh(rawToken, testSecret, testExpiry, testRefreshExpiry)

	assert.ErrorIs(t, err, services.ErrInvalidRefreshToken)
	refreshRepo.AssertExpectations(t)
}

func TestAuthService_Refresh_InvalidToken(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenRepository)
	svc := services.NewAuthService(userRepo, refreshRepo)

	rawToken := "unknown_token"
	tokenHash := services.HashRefreshToken(rawToken)

	refreshRepo.On("FindByTokenHash", tokenHash).Return(nil, errors.New("not found"))

	_, err := svc.Refresh(rawToken, testSecret, testExpiry, testRefreshExpiry)

	assert.ErrorIs(t, err, services.ErrInvalidRefreshToken)
	refreshRepo.AssertExpectations(t)
}

func TestAuthService_DeleteAccount(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenRepository)
	svc := services.NewAuthService(userRepo, refreshRepo)

	refreshRepo.On("DeleteAllByUserID", 15).Return(nil)
	userRepo.On("Delete", 15).Return(nil)

	err := svc.DeleteAccount(15)

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}
