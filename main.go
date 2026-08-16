package main

import (
  "fmt"
  // "time"
  "net/http"
  // "ollama-gui/internal/api"
)

func main() {
	fmt.Println("Server Startup")

    // rawURL := "http://localhost:11434"
    // rawToken := "secret-token"

    // client, err := api.NewClient(&rawURL, &rawToken, 5 * time.Second)
    // if err != nil {
    //     fmt.Println("error starting ollama client", err)
    // }
    
    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/", fileServer)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("error starting server", err)
    }

    fmt.Println("Server Exit")
}
