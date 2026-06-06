package mocks

import (
	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/Andrii-K-17/just-tasks/internal/repository"
	"github.com/stretchr/testify/mock"
)

// TaskRepository is a mock implementation of repository.TaskRepository.
type TaskRepository struct {
	mock.Mock
}

func (m *TaskRepository) FindAllByUser(userID int) ([]models.TaskResponse, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.([]models.TaskResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *TaskRepository) Create(userID int, p repository.TaskCreateParams) (*models.TaskResponse, error) {
	args := m.Called(userID, p)
	if v := args.Get(0); v != nil {
		return v.(*models.TaskResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *TaskRepository) GetOwnerID(taskID int) (int, error) {
	args := m.Called(taskID)
	return args.Int(0), args.Error(1)
}

func (m *TaskRepository) IsCollaborator(taskID, userID int) (bool, error) {
	args := m.Called(taskID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *TaskRepository) Update(taskID int, p repository.TaskUpdateParams) error {
	args := m.Called(taskID, p)
	return args.Error(0)
}

func (m *TaskRepository) Delete(userID, taskID int) (bool, error) {
	args := m.Called(userID, taskID)
	return args.Bool(0), args.Error(1)
}

func (m *TaskRepository) Reorder(userID int, ids []int) error {
	args := m.Called(userID, ids)
	return args.Error(0)
}

func (m *TaskRepository) FindCollaboratorByUsername(username string) (int, error) {
	args := m.Called(username)
	return args.Int(0), args.Error(1)
}

func (m *TaskRepository) AddCollaborator(taskID, userID int) (bool, error) {
	args := m.Called(taskID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *TaskRepository) RemoveCollaborator(taskID, collabID int) (bool, error) {
	args := m.Called(taskID, collabID)
	return args.Bool(0), args.Error(1)
}

func (m *TaskRepository) GetOwnerUsername(userID int) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}
