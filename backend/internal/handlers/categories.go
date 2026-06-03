package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Andrii-K-17/just-tasks/internal/middleware"
	"github.com/Andrii-K-17/just-tasks/internal/response"
	"github.com/Andrii-K-17/just-tasks/internal/services"
	"github.com/go-chi/chi/v5"
)

// CategoryHandler manages category HTTP endpoints.
type CategoryHandler struct {
	svc *services.CategoryService
}

// NewCategoryHandler initializes and returns a new CategoryHandler.
func NewCategoryHandler(svc *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// GetCategories retrieves all categories belonging to the authenticated user.
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	categories, err := h.svc.GetAll(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
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

	category, err := h.svc.Create(userID, name)
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

	deleted, err := h.svc.Delete(userID, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"deleted": deleted})
}
