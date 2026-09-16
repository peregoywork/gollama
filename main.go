package main

import (
	"fmt"
	"time"
	"log"
	"bytes"
	"os"
	"context"
	"net/http"
	"encoding/json"
	"path/filepath"
)

const (
	staticDir = "./public"
	promptsDir = "./prompts"
	host = "0.0.0.0"
	port = "8081"
	ollamaURL = "http://arcadia.home.arpa:11434" // "http://localhost:11434" 
)

func main() { 
	pacificLoc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Critical: failed to load timezone: %w", err)
	}

	ollamaClient := OllamaClient{
		http: &http.Client{ Timeout: time.Second * 120, },
		baseURL: ollamaURL,
	}

	server := Server{ ollama: &ollamaClient, }

	http.HandleFunc("/", middlewareLogging(pacificLoc, handleStaticFiles))
	http.HandleFunc("/chat", middlewareLogging(pacificLoc, server.handleOllamaChat))

	srvAddr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server starting: %s", srvAddr)
	err = http.ListenAndServe(srvAddr, nil)
    if err != nil {
		log.Fatalf("Critical: error starting server: %w", err)
    }

    fmt.Println("Server Exit")
}

type Server struct {
	ollama *OllamaClient
}

type OllamaClient struct {
	http 		*http.Client
	baseURL 	string
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


func getSystemPromptMessage() (ChatMessage, error) {
	filepath := filepath.Join(promptsDir, "system.md")
	sysPrompt, err := os.ReadFile(filepath)
	if err != nil {
		return ChatMessage{}, fmt.Errorf("failed to ready system prompt: %w", err)
	}

	msg := ChatMessage{
		Role: "system",
		Content: string(sysPrompt),
	}

	return msg, nil
}

func ensureSystemPrompt(messages []ChatMessage, sysMsg ChatMessage) []ChatMessage {
	if len(messages) == 0 || messages[0].Role != "system" {
		return append([]ChatMessage{sysMsg}, messages...)
	}
	return messages
}

func (c *OllamaClient) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	
	fullURL := c.baseURL + "/api/chat"
	reader := bytes.NewReader(data)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, reader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sent request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status: %d", resp.StatusCode)
	}

	var outResp ChatResponse
	if err = json.NewDecoder(resp.Body).Decode(&outResp); err != nil {
		return nil, fmt.Errorf("decoding error: %w", err)
	}

	return &outResp, nil
}


// Handlers

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
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	sysMsg, err := getSystemPromptMessage()
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return 
	}

	req.Messages = ensureSystemPrompt(req.Messages, sysMsg)

	resp, err := s.ollama.Chat(context.Background(), req)
	if err != nil {
		log.Printf("ollama error: %w", err)
		http.Error(w, "Failed to communicate with ollama server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
