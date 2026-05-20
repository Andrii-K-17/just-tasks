package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/andriik17/just-tasks/internal/middleware"
	"github.com/andriik17/just-tasks/internal/models"
	"github.com/andriik17/just-tasks/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// CategoryHandler manages category CRUD operations for tasks.
type CategoryHandler struct {
	db *sqlx.DB
}

// NewCategoryHandler initializes and returns a new CategoryHandler.
func NewCategoryHandler(db *sqlx.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// GetCategories retrieves all categories belonging to the authenticated user.
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var categories []models.Category
	err := h.db.Select(&categories, "SELECT id, name FROM categories WHERE user_id=$1 ORDER BY id ASC", userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	if categories == nil {
		categories = []models.Category{}
	}

	response.JSON(w, http.StatusOK, categories)
}

// createCategoryReq represents the payload for creating a new category.
type createCategoryReq struct {
	Name string `json:"name"`
}

// CreateCategory validates input and inserts a new category for the authenticated user.
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req createCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		response.Error(w, http.StatusUnprocessableEntity, "category name is required")
		return
	}

	var category models.Category
	err := h.db.QueryRowx(
		"INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id, name",
		userID, name,
	).StructScan(&category)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusCreated, category)
}

// DeleteCategory removes a category by ID if it belongs to the authenticated user.
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid category id")
		return
	}

	res, err := h.db.Exec("DELETE FROM categories WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	n, err := res.RowsAffected()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not determine affected rows")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"deleted": n > 0})
}
