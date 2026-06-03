package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/jmoiron/sqlx"
)

// ErrTaskNotFound is returned when a task does not exist.
var ErrTaskNotFound = errors.New("task not found")

// ErrForbidden is returned when a user lacks permission to perform an action.
var ErrForbidden = errors.New("forbidden")

// TaskService handles task CRUD, reordering, and collaborator management.
type TaskService struct {
	db *sqlx.DB
}

// NewTaskService initializes and returns a new TaskService.
func NewTaskService(db *sqlx.DB) *TaskService {
	return &TaskService{db: db}
}

// GetAll retrieves all owned and shared tasks for the given user.
func (s *TaskService) GetAll(userID int) ([]models.TaskResponse, error) {
	query := `
		SELECT t.id, t.user_id, u.username AS owner_name, t.task_text, t.priority,
			   TO_CHAR(t.deadline, 'YYYY-MM-DD') AS deadline,
			   t.is_completed, t.position, t.category_id, t.created_at
		FROM tasks t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id = $1 OR EXISTS (
			SELECT 1 FROM task_collaborators tc WHERE tc.task_id = t.id AND tc.user_id = $1
		)
		ORDER BY t.position ASC, t.id DESC
	`

	type dbTask struct {
		models.Task
		OwnerName string `db:"owner_name"`
	}
	var dbTasks []dbTask

	if err := s.db.Select(&dbTasks, query, userID); err != nil {
		return nil, err
	}

	var taskIDs []int
	for _, t := range dbTasks {
		taskIDs = append(taskIDs, t.ID)
	}

	collabsByTask := make(map[int][]models.TaskCollaborator)

	if len(taskIDs) > 0 {
		collabQuery, args, err := sqlx.In(`
			SELECT tc.task_id, u.id, u.username
			FROM task_collaborators tc
			JOIN users u ON tc.user_id = u.id
			WHERE tc.task_id IN (?)
		`, taskIDs)

		if err == nil {
			collabQuery = s.db.Rebind(collabQuery)
			type collabRow struct {
				TaskID   int    `db:"task_id"`
				ID       int    `db:"id"`
				Username string `db:"username"`
			}
			var collabRows []collabRow
			if err := s.db.Select(&collabRows, collabQuery, args...); err == nil {
				for _, row := range collabRows {
					collabsByTask[row.TaskID] = append(collabsByTask[row.TaskID], models.TaskCollaborator{
						ID:       row.ID,
						Username: row.Username,
					})
				}
			}
		}
	}

	var tasks []models.TaskResponse
	for _, dt := range dbTasks {
		collabs := collabsByTask[dt.ID]
		if collabs == nil {
			collabs = []models.TaskCollaborator{}
		}
		tasks = append(tasks, models.TaskResponse{
			Task:          dt.Task,
			OwnerName:     dt.OwnerName,
			Collaborators: collabs,
		})
	}

	if tasks == nil {
		tasks = []models.TaskResponse{}
	}

	return tasks, nil
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

	var task models.Task
	err := s.db.QueryRowx(
		`INSERT INTO tasks (user_id, task_text, priority, deadline, category_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, task_text, priority,
				   TO_CHAR(deadline, 'YYYY-MM-DD') AS deadline,
				   is_completed, position, category_id, created_at`,
		userID, strings.TrimSpace(p.TaskText), p.Priority, p.Deadline, p.CategoryID,
	).StructScan(&task)
	if err != nil {
		return nil, err
	}

	var ownerName string
	_ = s.db.Get(&ownerName, "SELECT username FROM users WHERE id=$1", userID)

	return &models.TaskResponse{
		Task:          task,
		OwnerName:     ownerName,
		Collaborators: []models.TaskCollaborator{},
	}, nil
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
	var taskOwner int
	if err := s.db.Get(&taskOwner, "SELECT user_id FROM tasks WHERE id=$1", taskID); err != nil {
		return ErrTaskNotFound
	}

	isOwner := taskOwner == userID
	isCollaborator := false

	if !isOwner {
		var exists bool
		if err := s.db.Get(&exists,
			"SELECT EXISTS(SELECT 1 FROM task_collaborators WHERE task_id=$1 AND user_id=$2)",
			taskID, userID,
		); err == nil {
			isCollaborator = exists
		}
	}

	if !isOwner && !isCollaborator {
		return ErrForbidden
	}

	sets := make([]string, 0)
	args := make([]any, 0)
	idx := 1

	if p.TaskText != nil {
		sets = append(sets, fmt.Sprintf("task_text=$%d", idx))
		args = append(args, *p.TaskText)
		idx++
	}
	if p.IsCompleted != nil {
		sets = append(sets, fmt.Sprintf("is_completed=$%d", idx))
		args = append(args, *p.IsCompleted)
		idx++
	}
	if p.Priority != nil {
		v := *p.Priority
		if v < 1 || v > 3 {
			v = 2
		}
		sets = append(sets, fmt.Sprintf("priority=$%d", idx))
		args = append(args, v)
		idx++
	}
	if p.DeadlineNull {
		sets = append(sets, fmt.Sprintf("deadline=$%d", idx))
		args = append(args, nil)
		idx++
	} else if p.Deadline != nil {
		sets = append(sets, fmt.Sprintf("deadline=$%d", idx))
		args = append(args, *p.Deadline)
		idx++
	}
	if p.CategoryNull {
		sets = append(sets, fmt.Sprintf("category_id=$%d", idx))
		args = append(args, nil)
		idx++
	} else if p.CategoryID != nil {
		sets = append(sets, fmt.Sprintf("category_id=$%d", idx))
		args = append(args, *p.CategoryID)
		idx++
	}

	if len(sets) == 0 {
		return errors.New("no fields to update")
	}

	args = append(args, taskID)
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id=$%d", strings.Join(sets, ", "), idx)

	res, err := s.db.Exec(query, args...)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		return errors.New("could not update task")
	}

	return nil
}

// Delete removes a task by ID if the user is its owner.
func (s *TaskService) Delete(userID, taskID int) error {
	res, err := s.db.Exec(
		"DELETE FROM tasks WHERE id=$1 AND user_id=$2",
		taskID, userID,
	)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrForbidden
	}

	return nil
}

// Reorder updates each task's position within a transaction.
func (s *TaskService) Reorder(userID int, ids []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		UPDATE tasks
		SET position = $1
		WHERE id = $2
			AND (user_id = $3 OR id IN (
				SELECT task_id FROM task_collaborators WHERE user_id = $3
			))
	`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for pos, id := range ids {
		if _, err := stmt.Exec(pos, id, userID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// AddCollaborator adds a user by username to a task owned by userID.
func (s *TaskService) AddCollaborator(ownerID, taskID int, username string) (bool, error) {
	var taskOwnerID int
	if err := s.db.Get(&taskOwnerID,
		"SELECT user_id FROM tasks WHERE id=$1",
		taskID,
	); err != nil || taskOwnerID != ownerID {
		return false, ErrForbidden
	}

	var collabID int
	if err := s.db.Get(&collabID,
		"SELECT id FROM users WHERE username=$1",
		username,
	); err != nil {
		return false, ErrTaskNotFound
	}

	if collabID == ownerID {
		return false, errors.New("сannot add yourself as a collaborator")
	}

	res, err := s.db.Exec(
		"INSERT INTO task_collaborators (task_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		taskID, collabID,
	)
	if err != nil {
		return false, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

// RemoveCollaborator removes a collaborator from a task owned by ownerID.
func (s *TaskService) RemoveCollaborator(ownerID, taskID, collabID int) (bool, error) {
	var taskOwnerID int
	if err := s.db.Get(&taskOwnerID,
		"SELECT user_id FROM tasks WHERE id=$1",
		taskID,
	); err != nil || taskOwnerID != ownerID {
		return false, ErrForbidden
	}

	res, err := s.db.Exec(
		"DELETE FROM task_collaborators WHERE task_id=$1 AND user_id=$2",
		taskID, collabID,
	)
	if err != nil {
		return false, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}
