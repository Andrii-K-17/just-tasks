package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ErrAINotConfigured is returned when no Groq API key is set.
var ErrAINotConfigured = errors.New("AI service is not configured")

// AITask represents a single generated task with metadata.
type AITask struct {
	Text     string `json:"text"`
	Deadline string `json:"deadline"`
	Priority string `json:"priority"`
}

// AIGenerateResult represents the full AI generation result.
type AIGenerateResult struct {
	Category string   `json:"category"`
	Tasks    []AITask `json:"tasks"`
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

// AIService handles AI-powered task generation via the Groq API.
type AIService struct {
	groqAPIKey string
	httpClient *http.Client
}

// NewAIService initializes and returns a new AIService.
func NewAIService(groqAPIKey string) *AIService {
	return &AIService{
		groqAPIKey: groqAPIKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GenerateTasks sends the user prompt to Groq and returns a structured task plan.
func (s *AIService) GenerateTasks(ctx context.Context, text, today string) (*AIGenerateResult, error) {
	if s.groqAPIKey == "" {
		return nil, ErrAINotConfigured
	}

	systemPrompt := fmt.Sprintf(
		`You are a task planning assistant.
		 Given a goal, generate a category name and a list of actionable subtasks with deadlines and priorities.
		 Today is %s. Respond ONLY with valid JSON, no markdown, no explanation.
		 Sort tasks by deadline ascending.
		 Use this exact format:
		 {"category":"string","tasks":[{"text":"string","deadline":"YYYY-MM-DD","priority":"high|medium|low"}]}`,
		today,
	)

	body, err := json.Marshal(groqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []groqMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: text},
		},
		MaxTokens:   1024,
		Temperature: 0.7,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.groqAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach AI service: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var groqResp groqResponse
	if err := json.Unmarshal(respBytes, &groqResp); err != nil || len(groqResp.Choices) == 0 {
		return nil, errors.New("invalid AI response structure")
	}

	content := strings.TrimSpace(groqResp.Choices[0].Message.Content)

	if strings.HasPrefix(content, "```") {
		lines := strings.Split(content, "\n")
		if len(lines) > 2 {
			content = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	var result AIGenerateResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, errors.New("failed to parse AI output")
	}

	return &result, nil
}
