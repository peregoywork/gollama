package main

import (
	"fmt"
	"time"
	"log"
	"net/http"

	// "ollama-gui/internal/api"
)

const (
	staticDir = "./public"
	host = "0.0.0.0"
	port = "8080"
)

func main() {
	url := fmt.Sprintf("%s:%s", host, port)
	pacificLoc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Critical: failed to load timezone: %w", err)
	}

	http.HandleFunc("/", middlewareLogging(pacificLoc, handleStaticFiles))

	log.Printf("Server starting: http://%s", url)
	err = http.ListenAndServe(url, nil)
    if err != nil {
		log.Fatalf("Critical: error starting server: %w", err)
    }

    fmt.Println("Server Exit")
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
