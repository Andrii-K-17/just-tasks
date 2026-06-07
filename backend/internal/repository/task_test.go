package repository_test

import (
	"testing"

	"github.com/Andrii-K-17/just-tasks/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_CreateAndList(t *testing.T) {
	db := newTestDB(t)
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	user, err := userRepo.Create("user1", "hash")
	require.NoError(t, err)

	task, err := taskRepo.Create(user.ID, repository.TaskCreateParams{
		TaskText: "test",
		Priority: 2,
	})
	require.NoError(t, err)
	assert.Equal(t, "test", task.TaskText)

	tasks, err := taskRepo.FindAllByUser(user.ID)
	require.NoError(t, err)
	assert.Len(t, tasks, 1)
}

func TestTaskRepository_Update(t *testing.T) {
	db := newTestDB(t)
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	user, err := userRepo.Create("user1", "hash")
	require.NoError(t, err)

	task, err := taskRepo.Create(user.ID, repository.TaskCreateParams{
		TaskText: "original",
		Priority: 1,
	})
	require.NoError(t, err)

	newText := "updated"
	err = taskRepo.Update(task.ID, repository.TaskUpdateParams{TaskText: &newText})
	require.NoError(t, err)

	tasks, err := taskRepo.FindAllByUser(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated", tasks[0].TaskText)
}

func TestTaskRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	user, err := userRepo.Create("user1", "hash")
	require.NoError(t, err)

	task, err := taskRepo.Create(user.ID, repository.TaskCreateParams{
		TaskText: "to delete",
		Priority: 2,
	})
	require.NoError(t, err)

	deleted, err := taskRepo.Delete(user.ID, task.ID)
	require.NoError(t, err)
	assert.True(t, deleted)

	tasks, err := taskRepo.FindAllByUser(user.ID)
	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskRepository_Collaborator(t *testing.T) {
	db := newTestDB(t)
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	owner, err := userRepo.Create("owner", "hash")
	require.NoError(t, err)
	collab, err := userRepo.Create("collab", "hash")
	require.NoError(t, err)

	task, err := taskRepo.Create(owner.ID, repository.TaskCreateParams{
		TaskText: "shared task",
		Priority: 2,
	})
	require.NoError(t, err)

	added, err := taskRepo.AddCollaborator(task.ID, collab.ID)
	require.NoError(t, err)
	assert.True(t, added)

	isCollab, err := taskRepo.IsCollaborator(task.ID, collab.ID)
	require.NoError(t, err)
	assert.True(t, isCollab)

	removed, err := taskRepo.RemoveCollaborator(task.ID, collab.ID)
	require.NoError(t, err)
	assert.True(t, removed)
}
