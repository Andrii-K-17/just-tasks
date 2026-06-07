package repository_test

import (
	"testing"

	"github.com/Andrii-K-17/just-tasks/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_CreateAndList(t *testing.T) {
	db := newTestDB(t)
	userRepo := repository.NewUserRepository(db)
	catRepo := repository.NewCategoryRepository(db)

	user, err := userRepo.Create("user1", "hash")
	require.NoError(t, err)

	cat, err := catRepo.Create(user.ID, "work")
	require.NoError(t, err)
	assert.Equal(t, "work", cat.Name)

	all, err := catRepo.FindAllByUser(user.ID)
	require.NoError(t, err)
	assert.Len(t, all, 1)
}

func TestCategoryRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	userRepo := repository.NewUserRepository(db)
	catRepo := repository.NewCategoryRepository(db)

	user, err := userRepo.Create("user1", "hash")
	require.NoError(t, err)

	cat, err := catRepo.Create(user.ID, "personal")
	require.NoError(t, err)

	deleted, err := catRepo.Delete(user.ID, cat.ID)
	require.NoError(t, err)
	assert.True(t, deleted)

	all, err := catRepo.FindAllByUser(user.ID)
	require.NoError(t, err)
	assert.Empty(t, all)
}
