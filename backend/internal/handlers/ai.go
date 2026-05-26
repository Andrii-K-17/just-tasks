package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/andriik17/just-tasks/internal/response"
)

// AIHandler manages AI-powered task generation via the Groq API.
type AIHandler struct {
	groqAPIKey string
}

// NewAIHandler initializes and returns a new AIHandler.
func NewAIHandler(groqAPIKey string) *AIHandler {
	return &AIHandler{groqAPIKey: groqAPIKey}
}

// aiGenerateRequest represents the user prompt payload.
type aiGenerateRequest struct {
	Text string `json:"text"`
}

// aiTask represents a single generated task with metadata.
type aiTask struct {
	Text     string `json:"text"`
	Deadline string `json:"deadline"`
	Priority string `json:"priority"`
}

// aiGenerateResponse represents the full AI generation result.
type aiGenerateResponse struct {
	Category string   `json:"category"`
	Tasks    []aiTask `json:"tasks"`
}

// groqMessage represents a single message in the Groq chat format.
type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// groqRequest represents the request body sent to the Groq API.
type groqRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

// groqResponse represents the response received from the Groq API.
type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GenerateTasks calls the Groq API with the user's prompt and returns a structured task plan.
func (h *AIHandler) GenerateTasks(w http.ResponseWriter, r *http.Request) {
	if h.groqAPIKey == "" {
		response.Error(w, http.StatusServiceUnavailable, "AI service is not configured")
		return
	}

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
	systemPrompt := fmt.Sprintf(
		`You are a task planning assistant. Given a goal, generate a category name and a list of actionable subtasks with deadlines and priorities. Today is %s. Respond ONLY with valid JSON, no markdown, no explanation. Sort tasks by deadline ascending. Use this exact format:
{"category":"string","tasks":[{"text":"string","deadline":"YYYY-MM-DD","priority":"high|medium|low"}]}`,
		today,
	)

	body, err := json.Marshal(groqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []groqMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: req.Text},
		},
		MaxTokens:   1024,
		Temperature: 0.7,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	groqReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost,
		"https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	groqReq.Header.Set("Content-Type", "application/json")
	groqReq.Header.Set("Authorization", "Bearer "+h.groqAPIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(groqReq)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "failed to reach AI service")
		return
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to read AI response")
		return
	}

	var groqResp groqResponse
	if err := json.Unmarshal(respBytes, &groqResp); err != nil || len(groqResp.Choices) == 0 {
		response.Error(w, http.StatusInternalServerError, "invalid AI response structure")
		return
	}

	content := strings.TrimSpace(groqResp.Choices[0].Message.Content)

	if strings.HasPrefix(content, "```") {
		lines := strings.Split(content, "\n")
		if len(lines) > 2 {
			content = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	var result aiGenerateResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to parse AI output")
		return
	}

	response.JSON(w, http.StatusOK, result)
}
