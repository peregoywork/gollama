package api

import (
    "time"
)

// ChatRequest represents the payload for /api/chat
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   *bool         `json:"stream,omitempty"`
	Options  map[string]any `json:"options,omitempty"`
}

// ChatMessage represents an individual message in a chat history
type ChatMessage struct {
	Role    string `json:"role"`    // "system", "user", or "assistant"
	Content string `json:"content"`
}

// ChatResponse represents the server's JSON response for /api/chat
type ChatResponse struct {
	Model      string      `json:"model"`
	CreatedAt  time.Time   `json:"created_at"`
	Message    ChatMessage `json:"message"`
	Done       bool        `json:"done"`
	TotalDuration int64    `json:"total_duration"`
}

// GenerateRequest represents the payload for text completion (/api/generate)
type GenerateRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Stream  *bool          `json:"stream,omitempty"`
	Options map[string]any `json:"options,omitempty"`
}

// GenerateResponse represents the server's JSON response for /api/generate
type GenerateResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}
