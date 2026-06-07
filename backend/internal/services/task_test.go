package services_test

import (
	"errors"
	"testing"

	"github.com/Andrii-K-17/just-tasks/internal/mocks"
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/Andrii-K-17/just-tasks/internal/repository"
	"github.com/Andrii-K-17/just-tasks/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskService_GetAll(t *testing.T) {
	repo := new(mocks.TaskRepository)
	svc := services.NewTaskService(repo)

	expected := []models.TaskResponse{{Task: models.Task{ID: 1, TaskText: "test"}}}
	repo.On("FindAllByUser", 1).Return(expected, nil)

	tasks, err := svc.GetAll(1)

	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	repo.AssertExpectations(t)
}

func TestTaskService_Create_DefaultPriority(t *testing.T) {
	repo := new(mocks.TaskRepository)
	svc := services.NewTaskService(repo)

	expected := &models.TaskResponse{Task: models.Task{ID: 1, TaskText: "task", Priority: 2}}

	repo.On("Create", 1, repository.TaskCreateParams{
		TaskText: "task",
		Priority: 2,
	}).Return(expected, nil)

	task, err := svc.Create(1, services.CreatePayload{TaskText: "task", Priority: 0})

	require.NoError(t, err)
	assert.Equal(t, 2, task.Priority)
	repo.AssertExpectations(t)
}

func TestTaskService_Update_Forbidden(t *testing.T) {
	repo := new(mocks.TaskRepository)
	svc := services.NewTaskService(repo)

	repo.On("GetOwnerID", 10).Return(99, nil)
	repo.On("IsCollaborator", 10, 1).Return(false, nil)

	err := svc.Update(1, 10, services.UpdatePayload{})

	assert.ErrorIs(t, err, services.ErrForbidden)
	repo.AssertExpectations(t)
}

func TestTaskService_Update_OwnerCanEdit(t *testing.T) {
	repo := new(mocks.TaskRepository)
	svc := services.NewTaskService(repo)

	text := "updated"
	repo.On("GetOwnerID", 10).Return(1, nil)
	repo.On("Update", 10, repository.TaskUpdateParams{TaskText: &text}).Return(nil)

	err := svc.Update(1, 10, services.UpdatePayload{TaskText: &text})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTaskService_Update_CollaboratorCanEdit(t *testing.T) {
	repo := new(mocks.TaskRepository)
	svc := services.NewTaskService(repo)

	done := true
	repo.On("GetOwnerID", 10).Return(55, nil)
	repo.On("IsCollaborator", 10, 2).Return(true, nil)
	repo.On("Update", 10, repository.TaskUpdateParams{IsCompleted: &done}).Return(nil)

	err := svc.Update(2, 10, services.UpdatePayload{IsCompleted: &done})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTaskService_Delete_NotOwner(t *testing.T) {
	repo := new(mocks.TaskRepository)
	svc := services.NewTaskService(repo)

	repo.On("Delete", 1, 10).Return(false, nil)

	err := svc.Delete(1, 10)

	assert.ErrorIs(t, err, services.ErrForbidden)
	repo.AssertExpectations(t)
}

func TestTaskService_AddCollaborator_CannotAddSelf(t *testing.T) {
	repo := new(mocks.TaskRepository)
	svc := services.NewTaskService(repo)

	repo.On("GetOwnerID", 5).Return(1, nil)
	repo.On("FindCollaboratorByUsername", "collaborator").Return(1, nil)

	_, err := svc.AddCollaborator(1, 5, "collaborator")

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestTaskService_AddCollaborator_UserNotFound(t *testing.T) {
	repo := new(mocks.TaskRepository)
	svc := services.NewTaskService(repo)

	repo.On("GetOwnerID", 5).Return(1, nil)
	repo.On("FindCollaboratorByUsername", "ghost").Return(0, errors.New("not found"))

	_, err := svc.AddCollaborator(1, 5, "ghost")

	assert.ErrorIs(t, err, services.ErrTaskNotFound)
	repo.AssertExpectations(t)
}
