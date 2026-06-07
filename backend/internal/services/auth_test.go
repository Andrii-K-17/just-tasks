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

func TestAuthService_Register_Success(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := services.NewAuthService(repo)

	expected := &models.User{ID: 1, Username: "user1", CreatedAt: time.Now()}

	repo.On("ExistsByUsername", "user1").Return(false, nil)
	repo.On("Create", "user1", mock.MatchedBy(func(s string) bool {
		return strings.HasPrefix(s, "$2a$")
	})).Return(expected, nil)

	user, err := svc.Register("user1", "password123")

	require.NoError(t, err)
	assert.Equal(t, expected.Username, user.Username)
	repo.AssertExpectations(t)
}

func TestAuthService_register_UsernameTaken(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := services.NewAuthService(repo)

	repo.On("ExistsByUsername", "user1").Return(true, nil)

	_, err := svc.Register("user1", "password123")

	assert.ErrorIs(t, err, services.ErrUsernameTaken)
	repo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := services.NewAuthService(repo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	stored := &models.User{ID: 1, Username: "user1", PasswordHash: string(hash)}

	repo.On("FindByUsername", "user1").Return(stored, nil)

	user, err := svc.Login("user1", "password123")

	require.NoError(t, err)
	assert.Equal(t, stored.ID, user.ID)
	repo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := services.NewAuthService(repo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	stored := &models.User{ID: 1, Username: "user1", PasswordHash: string(hash)}

	repo.On("FindByUsername", "user1").Return(stored, nil)

	_, err := svc.Login("user1", "wrong")

	assert.ErrorIs(t, err, services.ErrInvalidCredentials)
	repo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := services.NewAuthService(repo)

	repo.On("FindByUsername", "ghost").Return(nil, errors.New("not found"))
	_, err := svc.Login("ghost", "password123")

	assert.ErrorIs(t, err, services.ErrInvalidCredentials)
	repo.AssertExpectations(t)
}

func TestAuthService_DeleteAccount(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := services.NewAuthService(repo)

	repo.On("Delete", 15).Return(nil)
	err := svc.DeleteAccount(15)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
