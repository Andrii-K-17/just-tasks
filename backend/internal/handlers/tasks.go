package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/andriik17/just-tasks/internal/middleware"
	"github.com/andriik17/just-tasks/internal/models"
	"github.com/andriik17/just-tasks/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// TaskHandler manages task-related operations and database access.
type TaskHandler struct {
	db *sqlx.DB
}

// NewTaskHandler initializes and returns a new TaskHandler.
func NewTaskHandler(db *sqlx.DB) *TaskHandler {
	return &TaskHandler{db: db}
}

// GetTasks retrieves all owned and shared tasks for the authenticated user.
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

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

	err := h.db.Select(&dbTasks, query, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
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
			collabQuery = h.db.Rebind(collabQuery)
			type CollabRow struct {
				TaskID   int    `db:"task_id"`
				ID       int    `db:"id"`
				Username string `db:"username"`
			}
			var collabRows []CollabRow
			if err := h.db.Select(&collabRows, collabQuery, args...); err == nil {
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

	response.JSON(w, http.StatusOK, tasks)
}

// createTaskRequest represents the payload for creating a new task.
type createTaskRequest struct {
	TaskText   string  `json:"task_text"`
	Priority   int     `json:"priority"`
	Deadline   *string `json:"deadline"`
	CategoryID *int    `json:"category_id"`
}

// CreateTask inserts a new task into the database and returns it.
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.TaskText) == "" {
		response.Error(w, http.StatusUnprocessableEntity, "the task text field is required")
		return
	}
	if req.Priority < 1 || req.Priority > 3 {
		req.Priority = 2
	}

	var deadline *string
	if req.Deadline != nil && *req.Deadline != "" {
		deadline = req.Deadline
	}

	var task models.Task
	err := h.db.QueryRowx(
		`INSERT INTO tasks (user_id, task_text, priority, deadline, category_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, task_text, priority,
				   TO_CHAR(deadline, 'YYYY-MM-DD') AS deadline,
				   is_completed, position, category_id, created_at`,
		userID, strings.TrimSpace(req.TaskText), req.Priority, deadline, req.CategoryID,
	).StructScan(&task)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	var ownerName string
	h.db.Get(&ownerName, "SELECT username FROM users WHERE id=$1", userID)

	res := models.TaskResponse{
		Task:          task,
		OwnerName:     ownerName,
		Collaborators: []models.TaskCollaborator{},
	}

	response.JSON(w, http.StatusCreated, res)
}

// UpdateTask modifies fields of an existing task belonging to the authenticated user.
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var taskOwner int
	err = h.db.Get(&taskOwner, "SELECT user_id FROM tasks WHERE id=$1", id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "task not found")
		return
	}

	isOwner := (taskOwner == userID)
	isCollaborator := false

	if !isOwner {
		var exists bool
		err := h.db.Get(&exists,
			"SELECT EXISTS(SELECT 1 FROM task_collaborators WHERE task_id=$1 AND user_id=$2)",
			id, userID,
		)
		if err == nil {
			isCollaborator = exists
		}
	}

	if !isOwner && !isCollaborator {
		response.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sets := make([]string, 0, len(body))
	args := make([]any, 0, len(body)+1)
	idx := 1

	if v, ok := body["task_text"]; ok {
		text := strings.TrimSpace(fmt.Sprintf("%v", v))
		if text == "" {
			response.Error(w, http.StatusUnprocessableEntity, "the task text field is required")
			return
		}
		sets = append(sets, fmt.Sprintf("task_text=$%d", idx))
		args = append(args, text)
		idx++
	}

	if v, ok := body["is_completed"]; ok {
		sets = append(sets, fmt.Sprintf("is_completed=$%d", idx))
		args = append(args, v.(bool))
		idx++
	}

	if v, ok := body["priority"]; ok {
		p := int(v.(float64))
		if p < 1 || p > 3 {
			p = 2
		}
		sets = append(sets, fmt.Sprintf("priority=$%d", idx))
		args = append(args, p)
		idx++
	}

	if v, ok := body["deadline"]; ok {
		sets = append(sets, fmt.Sprintf("deadline=$%d", idx))
		if v == nil || v == "" {
			args = append(args, nil)
		} else {
			args = append(args, v.(string))
		}
		idx++
	}

	if v, ok := body["category_id"]; ok {
		sets = append(sets, fmt.Sprintf("category_id=$%d", idx))
		if v == nil {
			args = append(args, nil)
		} else {
			args = append(args, int(v.(float64)))
		}
		idx++
	}

	if len(sets) == 0 {
		response.Error(w, http.StatusUnprocessableEntity, "no fields to update")
		return
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE tasks SET %s WHERE id=$%d",
		strings.Join(sets, ", "), idx,
	)

	res, err := h.db.Exec(query, args...)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		response.Error(w, http.StatusInternalServerError, "could not update task")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"updated": n > 0})
}

// DeleteTask removes a task by its ID for the authenticated user.
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	res, err := h.db.Exec("DELETE FROM tasks WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	n, err := res.RowsAffected()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not determine affected rows")
		return
	}
	if n == 0 {
		response.Error(w, http.StatusForbidden, "forbidden or task not found")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"deleted": n > 0})
}

// reorderRequest represents the payload for reordering tasks.
type reorderRequest struct {
	IDs []int `json:"ids"`
}

// ReorderTasks updates the position of tasks based on the provided order.
func (h *TaskHandler) ReorderTasks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req reorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.IDs) == 0 {
		response.Error(w, http.StatusUnprocessableEntity, "ids array required")
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE tasks SET position=$1 WHERE id=$2 AND (user_id=$3 OR id IN (SELECT task_id FROM task_collaborators WHERE user_id=$3))")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer stmt.Close()

	for pos, id := range req.IDs {
		if _, err := stmt.Exec(pos, id, userID); err != nil {
			response.Error(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// AddCollaborator adds a user to a task by their username.
func (h *TaskHandler) AddCollaborator(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" {
		response.Error(w, http.StatusBadRequest, "username required")
		return
	}

	var ownerID int
	err = h.db.Get(&ownerID, "SELECT user_id FROM tasks WHERE id=$1", taskID)
	if err != nil || ownerID != userID {
		response.Error(w, http.StatusForbidden, "only owner can add collaborators")
		return
	}

	var collabID int
	err = h.db.Get(&collabID, "SELECT id FROM users WHERE username=$1", req.Username)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	if collabID == userID {
		response.Error(w, http.StatusBadRequest, "сannot add yourself as a collaborator")
		return
	}

	res, err := h.db.Exec("INSERT INTO task_collaborators (task_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", taskID, collabID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	n, err := res.RowsAffected()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not determine affected rows")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"added": n > 0})
}

// RemoveCollaborator removes a user from a task.
func (h *TaskHandler) RemoveCollaborator(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	collabIDStr := chi.URLParam(r, "collabId")
	collabID, err := strconv.Atoi(collabIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid collaborator id")
		return
	}

	var ownerID int
	err = h.db.Get(&ownerID, "SELECT user_id FROM tasks WHERE id=$1", taskID)
	if err != nil || ownerID != userID {
		response.Error(w, http.StatusForbidden, "only owner can remove collaborators")
		return
	}

	res, err := h.db.Exec("DELETE FROM task_collaborators WHERE task_id=$1 AND user_id=$2", taskID, collabID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	n, err := res.RowsAffected()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not determine affected rows")
		return
	}
	response.JSON(w, http.StatusOK, map[string]bool{"removed": n > 0})
}
