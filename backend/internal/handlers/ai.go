package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Andrii-K-17/just-tasks/internal/response"
	"github.com/Andrii-K-17/just-tasks/internal/services"
)

// AIHandler manages AI-powered task generation endpoints.
type AIHandler struct {
	svc *services.AIService
}

// NewAIHandler initializes and returns a new AIHandler.
func NewAIHandler(svc *services.AIService) *AIHandler {
	return &AIHandler{svc: svc}
}

// aiGenerateRequest represents the user prompt payload.
type aiGenerateRequest struct {
	Text string `json:"text"`
}

// GenerateTasks calls the Groq API with the user's prompt and returns a structured task plan.
func (h *AIHandler) GenerateTasks(w http.ResponseWriter, r *http.Request) {
	var req aiGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		response.Error(w, http.StatusUnprocessableEntity, "text is required")
		return
	}

	today := time.Now().Format("2006-01-02")

	result, err := h.svc.GenerateTasks(r.Context(), req.Text, today)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAINotConfigured):
			response.Error(w, http.StatusServiceUnavailable, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, result)
}
