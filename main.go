package main

import (
	"fmt"
	"time"
	"log"
	"net/http"
	"encoding/json"

	"ollama-gui/internal/api"
)

const (
	staticDir = "./public"
	host = "0.0.0.0"
	port = "8080"
	ollamaURL = "localhost:11434"
)

func main() {
	url := fmt.Sprintf("%s:%s", host, port)
	pacificLoc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Critical: failed to load timezone: %w", err)
	}

	ollamaClient, err := api.NewClient(&url, nil)
	if err != nil {
		log.Fatalf("Critical: failed to create api client: %w", err)
	}

	http.HandleFunc("/", middlewareLogging(pacificLoc, handleStaticFiles))
	http.HandleFunc("/chat", middlewareLogging(pacificLoc, handleOllamaChat))

	log.Printf("Server starting: http://%s", url)
	err = http.ListenAndServe(url, nil)
    if err != nil {
		log.Fatalf("Critical: error starting server: %w", err)
    }

    fmt.Println("Server Exit")
}



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


func middlewareLogging(loc *time.Location, next http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		start := time.Now().In(loc)
		next(w, r)
		log.Printf(
			"%s %s | %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	}
}

func handleStaticFiles(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir(staticDir))
	fs.ServeHTTP(w, r)
}

// Server method: Handles user HTTP traffic
func (s *Server) handleOllamaChat(w http.ResponseWriter, r *http.Request) {
    // 1. Read user JSON from r.Body
    // 2. Call Ollama via s.ollama.Chat(r.Context(), userReq)
    // 3. Write response to w
}
