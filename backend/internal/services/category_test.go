package services_test

import (
	"testing"

	"github.com/Andrii-K-17/just-tasks/internal/mocks"
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/Andrii-K-17/just-tasks/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryService_GetAll(t *testing.T) {
	repo := new(mocks.CategoryRepository)
	svc := services.NewCategoryService(repo)

	expected := []models.Category{{ID: 1, Name: "work"}}
	repo.On("FindAllByUser", 1).Return(expected, nil)

	categories, err := svc.GetAll(1)

	require.NoError(t, err)
	assert.Len(t, categories, 1)
	repo.AssertExpectations(t)
}

func TestCategoryService_Create(t *testing.T) {
	repo := new(mocks.CategoryRepository)
	svc := services.NewCategoryService(repo)

	expected := &models.Category{ID: 2, Name: "personal"}
	repo.On("Create", 1, "personal").Return(expected, nil)

	category, err := svc.Create(1, "personal")

	require.NoError(t, err)
	assert.Equal(t, "personal", category.Name)
	repo.AssertExpectations(t)
}

func TestCategoryService_Delete(t *testing.T) {
	repo := new(mocks.CategoryRepository)
	svc := services.NewCategoryService(repo)

	repo.On("Delete", 1, 3).Return(true, nil)

	deleted, err := svc.Delete(1, 3)

	require.NoError(t, err)
	assert.True(t, deleted)
	repo.AssertExpectations(t)
}
