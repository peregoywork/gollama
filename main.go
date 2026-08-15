package main

import (
  "fmt"
  "time"
 
  "ollama-gui/internal/api"
)

func main() {
	fmt.Println("Server Startup")

    rawURL := "http://localhost:11434"
    rawToken := "secret-token"

    client, err := api.NewClient(&rawURL, &rawToken, 5 * time.Second)
    if err != nil {
        fmt.Println("error starting ollama client", err)
    }
    
    client.DoNothing()

    fmt.Println("Server Exit")
}
