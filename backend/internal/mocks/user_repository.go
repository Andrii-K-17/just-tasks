package mocks

import (
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/stretchr/testify/mock"
)

// UserRepository is a mock implementation of repository.UserRepository.
type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) ExistsByUsername(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepository) Create(username, passwordHash string) (*models.User, error) {
	args := m.Called(username, passwordHash)
	if v := args.Get(0); v != nil {
		return v.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if v := args.Get(0); v != nil {
		return v.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) FindByID(id int) (*models.User, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
