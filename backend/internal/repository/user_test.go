package repository_test

import (
	"testing"

	"github.com/Andrii-K-17/just-tasks/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateAndFind(t *testing.T) {
	db := newTestDB(t)
	repo := repository.NewUserRepository(db)

	user, err := repo.Create("user1", "hash123")
	require.NoError(t, err)
	assert.Equal(t, "user1", user.Username)
	assert.Positive(t, user.ID)

	found, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepository_ExistsByUsername(t *testing.T) {
	db := newTestDB(t)
	repo := repository.NewUserRepository(db)

	_, err := repo.Create("user1", "hash")
	require.NoError(t, err)

	exists, err := repo.ExistsByUsername("user1")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByUsername("nobody")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := repository.NewUserRepository(db)

	user, err := repo.Create("user1", "hash")
	require.NoError(t, err)

	err = repo.Delete(user.ID)
	require.NoError(t, err)

	_, err = repo.FindByID(user.ID)
	assert.Error(t, err)
}
