package main

import (
	"fmt"
	// "time"
	"net/http"

	// "ollama-gui/internal/api"
)

const (
	staticDir = "./public"
)

func main() {
	fmt.Println("Server Startup")
    // rawToken := "secret-token"

	http.HandleFunc("/", handleStaticFiles)

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("error starting server", err)
    }

    fmt.Println("Server Exit")
}


func handleStaticFiles(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir(staticDir))
	fs.ServeHTTP(w, r)
}
