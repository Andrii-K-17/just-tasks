package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Andrii-K-17/just-tasks/internal/middleware"
	"github.com/Andrii-K-17/just-tasks/internal/response"
	"github.com/Andrii-K-17/just-tasks/internal/services"
	"github.com/go-chi/chi/v5"
)

// TaskHandler manages task-related HTTP endpoints.
type TaskHandler struct {
	svc *services.TaskService
}

// NewTaskHandler initializes and returns a new TaskHandler.
func NewTaskHandler(svc *services.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// GetTasks retrieves all owned and shared tasks for the authenticated user.
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	tasks, err := h.svc.GetAll(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
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

	var deadline *string
	if req.Deadline != nil && *req.Deadline != "" {
		deadline = req.Deadline
	}

	task, err := h.svc.Create(userID, services.CreatePayload{
		TaskText:   req.TaskText,
		Priority:   req.Priority,
		Deadline:   deadline,
		CategoryID: req.CategoryID,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusCreated, task)
}

// UpdateTask modifies fields of an existing task belonging to the authenticated user.
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p := services.UpdatePayload{}

	if v, ok := body["task_text"]; ok {
		text := strings.TrimSpace(fmt.Sprintf("%v", v))
		if text == "" {
			response.Error(w, http.StatusUnprocessableEntity, "the task text field is required")
			return
		}
		p.TaskText = &text
	}
	if v, ok := body["is_completed"]; ok {
		b := v.(bool)
		p.IsCompleted = &b
	}
	if v, ok := body["priority"]; ok {
		pri := int(v.(float64))
		p.Priority = &pri
	}
	if v, ok := body["deadline"]; ok {
		if v == nil || v == "" {
			p.DeadlineNull = true
		} else {
			s := v.(string)
			p.Deadline = &s
		}
	}
	if v, ok := body["category_id"]; ok {
		if v == nil {
			p.CategoryNull = true
		} else {
			id := int(v.(float64))
			p.CategoryID = &id
		}
	}

	if err := h.svc.Update(userID, id, p); err != nil {
		switch {
		case errors.Is(err, services.ErrTaskNotFound):
			response.Error(w, http.StatusNotFound, "task not found")
		case errors.Is(err, services.ErrForbidden):
			response.Error(w, http.StatusForbidden, "forbidden")
		default:
			response.Error(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"updated": true})
}

// DeleteTask removes a task by its ID for the authenticated user.
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	if err := h.svc.Delete(userID, id); err != nil {
		if errors.Is(err, services.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "forbidden or task not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true})
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

	if err := h.svc.Reorder(userID, req.IDs); err != nil {
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

	added, err := h.svc.AddCollaborator(userID, taskID, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			response.Error(w, http.StatusForbidden, "only owner can add collaborators")
		case errors.Is(err, services.ErrTaskNotFound):
			response.Error(w, http.StatusNotFound, "user not found")
		default:
			response.Error(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"added": added})
}

// RemoveCollaborator removes a user from a task.
func (h *TaskHandler) RemoveCollaborator(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	collabID, err := strconv.Atoi(chi.URLParam(r, "collabId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid collaborator id")
		return
	}

	removed, err := h.svc.RemoveCollaborator(userID, taskID, collabID)
	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "only owner can remove collaborators")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"removed": removed})
}
