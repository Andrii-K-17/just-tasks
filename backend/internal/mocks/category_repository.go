package mocks

import (
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/stretchr/testify/mock"
)

// CategoryRepository is a mock implementation of repository.CategoryRepository.
type CategoryRepository struct {
	mock.Mock
}

func (m *CategoryRepository) FindAllByUser(userID int) ([]models.Category, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.([]models.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryRepository) Create(userID int, name string) (*models.Category, error) {
	args := m.Called(userID, name)
	if v := args.Get(0); v != nil {
		return v.(*models.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryRepository) Delete(userID, categoryID int) (bool, error) {
	args := m.Called(userID, categoryID)
	return args.Bool(0), args.Error(1)
}
