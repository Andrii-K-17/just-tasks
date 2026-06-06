package services

import (
	"errors"
	"strings"

	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/Andrii-K-17/just-tasks/internal/repository"
)

// ErrTaskNotFound is returned when a task does not exist.
var ErrTaskNotFound = errors.New("task not found")

// ErrForbidden is returned when a user lacks permission to perform an action.
var ErrForbidden = errors.New("forbidden")

// TaskService handles task CRUD, reordering, and collaborator management.
type TaskService struct {
	repo repository.TaskRepository
}

// NewTaskService initializes and returns a new TaskService.
func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// GetAll retrieves all owned and shared tasks for the given user.
func (s *TaskService) GetAll(userID int) ([]models.TaskResponse, error) {
	return s.repo.FindAllByUser(userID)
}

// CreatePayload holds the fields required to create a new task.
type CreatePayload struct {
	TaskText   string
	Priority   int
	Deadline   *string
	CategoryID *int
}

// Create inserts a new task and returns the created record with owner name.
func (s *TaskService) Create(userID int, p CreatePayload) (*models.TaskResponse, error) {
	if p.Priority < 1 || p.Priority > 3 {
		p.Priority = 2
	}

	return s.repo.Create(userID, repository.TaskCreateParams{
		TaskText:   strings.TrimSpace(p.TaskText),
		Priority:   p.Priority,
		Deadline:   p.Deadline,
		CategoryID: p.CategoryID,
	})
}

// UpdatePayload holds the optional fields that can be patched on a task.
type UpdatePayload struct {
	TaskText     *string
	IsCompleted  *bool
	Priority     *int
	Deadline     *string
	DeadlineNull bool
	CategoryID   *int
	CategoryNull bool
}

// Update applies the given patch to a task after verifying ownership or collaboration.
func (s *TaskService) Update(userID, taskID int, p UpdatePayload) error {
	ownerID, err := s.repo.GetOwnerID(taskID)
	if err != nil {
		return ErrTaskNotFound
	}

	isOwner := ownerID == userID
	if !isOwner {
		isCollab, err := s.repo.IsCollaborator(taskID, userID)
		if err != nil || !isCollab {
			return ErrForbidden
		}
	}

	return s.repo.Update(taskID, repository.TaskUpdateParams{
		TaskText:     p.TaskText,
		IsCompleted:  p.IsCompleted,
		Priority:     p.Priority,
		Deadline:     p.Deadline,
		DeadlineNull: p.DeadlineNull,
		CategoryID:   p.CategoryID,
		CategoryNull: p.CategoryNull,
	})
}

// Delete removes a task by ID if the user is its owner.
func (s *TaskService) Delete(userID, taskID int) error {
	deleted, err := s.repo.Delete(userID, taskID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrForbidden
	}
	return nil
}

// Reorder updates each task's position within a transaction.
func (s *TaskService) Reorder(userID int, ids []int) error {
	return s.repo.Reorder(userID, ids)
}

// AddCollaborator adds a user by username to a task owned by userID.
func (s *TaskService) AddCollaborator(ownerID, taskID int, username string) (bool, error) {
	taskOwnerID, err := s.repo.GetOwnerID(taskID)
	if err != nil || taskOwnerID != ownerID {
		return false, ErrForbidden
	}

	collabID, err := s.repo.FindCollaboratorByUsername(username)
	if err != nil {
		return false, ErrTaskNotFound
	}

	if collabID == ownerID {
		return false, errors.New("сannot add yourself as a collaborator")
	}

	return s.repo.AddCollaborator(taskID, collabID)
}

// RemoveCollaborator removes a collaborator from a task owned by ownerID.
func (s *TaskService) RemoveCollaborator(ownerID, taskID, collabID int) (bool, error) {
	taskOwnerID, err := s.repo.GetOwnerID(taskID)
	if err != nil || taskOwnerID != ownerID {
		return false, ErrForbidden
	}

	return s.repo.RemoveCollaborator(taskID, collabID)
}
