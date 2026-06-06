package repository

import (
	"fmt"
	"strings"

	"github.com/Andrii-K-17/just-tasks/internal/models"
	"github.com/jmoiron/sqlx"
)

// TaskRepository defines the persistence interface for task operations.
type TaskRepository interface {
	FindAllByUser(userID int) ([]models.TaskResponse, error)
	Create(userID int, p TaskCreateParams) (*models.TaskResponse, error)
	GetOwnerID(taskID int) (int, error)
	IsCollaborator(taskID, userID int) (bool, error)
	Update(taskID int, p TaskUpdateParams) error
	Delete(userID, taskID int) (bool, error)
	Reorder(userID int, ids []int) error
	FindCollaboratorByUsername(username string) (int, error)
	AddCollaborator(taskID, userID int) (bool, error)
	RemoveCollaborator(taskID, collabID int) (bool, error)
	GetOwnerUsername(userID int) (string, error)
}

// TaskCreateParams holds the fields required to insert a new task.
type TaskCreateParams struct {
	TaskText   string
	Priority   int
	Deadline   *string
	CategoryID *int
}

// TaskUpdateParams holds the optional fields that can be patched on a task.
type TaskUpdateParams struct {
	TaskText     *string
	IsCompleted  *bool
	Priority     *int
	Deadline     *string
	DeadlineNull bool
	CategoryID   *int
	CategoryNull bool
}

// pgTaskRepository is a PostgreSQL-backed implementation of TaskRepository.
type pgTaskRepository struct {
	db *sqlx.DB
}

// NewTaskRepository initializes and returns a new pgTaskRepository.
func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return &pgTaskRepository{db: db}
}

// FindAllByUser retrieves all owned and shared tasks for the given user.
func (r *pgTaskRepository) FindAllByUser(userID int) ([]models.TaskResponse, error) {
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

	if err := r.db.Select(&dbTasks, query, userID); err != nil {
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
			collabQuery = r.db.Rebind(collabQuery)
			type collabRow struct {
				TaskID   int    `db:"task_id"`
				ID       int    `db:"id"`
				Username string `db:"username"`
			}
			var collabRows []collabRow
			if err := r.db.Select(&collabRows, collabQuery, args...); err == nil {
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

// Create inserts a new task and returns the created record with owner name.
func (r *pgTaskRepository) Create(userID int, p TaskCreateParams) (*models.TaskResponse, error) {
	var task models.Task
	err := r.db.QueryRowx(
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
	_ = r.db.Get(&ownerName, "SELECT username FROM users WHERE id=$1", userID)

	return &models.TaskResponse{
		Task:          task,
		OwnerName:     ownerName,
		Collaborators: []models.TaskCollaborator{},
	}, nil
}

// GetOwnerID returns the user_id of the task owner.
func (r *pgTaskRepository) GetOwnerID(taskID int) (int, error) {
	var ownerID int
	err := r.db.Get(&ownerID, "SELECT user_id FROM tasks WHERE id=$1", taskID)
	return ownerID, err
}

// IsCollaborator checks whether a user is a collaborator on the given task.
func (r *pgTaskRepository) IsCollaborator(taskID, userID int) (bool, error) {
	var exists bool
	err := r.db.Get(&exists,
		"SELECT EXISTS(SELECT 1 FROM task_collaborators WHERE task_id=$1 AND user_id=$2)",
		taskID, userID,
	)
	return exists, err
}

// Update applies a partial patch to a task record.
func (r *pgTaskRepository) Update(taskID int, p TaskUpdateParams) error {
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
		return fmt.Errorf("no fields to update")
	}

	args = append(args, taskID)
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id=$%d", strings.Join(sets, ", "), idx)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		return fmt.Errorf("could not update task")
	}

	return nil
}

// Delete removes a task by ID if the user is its owner.
func (r *pgTaskRepository) Delete(userID, taskID int) (bool, error) {
	res, err := r.db.Exec(
		"DELETE FROM tasks WHERE id=$1 AND user_id=$2",
		taskID, userID,
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

// Reorder updates each task's position within a transaction.
func (r *pgTaskRepository) Reorder(userID int, ids []int) error {
	tx, err := r.db.Begin()
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

// FindCollaboratorByUsername resolves a username to a user ID.
func (r *pgTaskRepository) FindCollaboratorByUsername(username string) (int, error) {
	var id int
	err := r.db.Get(&id, "SELECT id FROM users WHERE username=$1", username)
	return id, err
}

// AddCollaborator inserts a collaborator row, ignoring conflicts.
func (r *pgTaskRepository) AddCollaborator(taskID, userID int) (bool, error) {
	res, err := r.db.Exec(
		"INSERT INTO task_collaborators (task_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		taskID, userID,
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

// RemoveCollaborator deletes a collaborator row by task and user ID.
func (r *pgTaskRepository) RemoveCollaborator(taskID, collabID int) (bool, error) {
	res, err := r.db.Exec(
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

// GetOwnerUsername retrieves the username of the given user ID.
func (r *pgTaskRepository) GetOwnerUsername(userID int) (string, error) {
	var username string
	err := r.db.Get(&username, "SELECT username FROM users WHERE id=$1", userID)
	return username, err
}
