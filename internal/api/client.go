package api

import (
    "fmt"
    "time"
    "net/http"
	// "encoding/json"
)

type Client struct {
   http         *http.Client
   url          string
   authToken    string
}

func NewClient (baseURL, token *string, timeout time.Duration) (*Client, error) {
    resolvedURL := "http://localhost:11434"
    resolvedToken := ""

    if baseURL != nil && *baseURL != "" {
        resolvedURL = *baseURL
    }

    if token != nil {
        resolvedToken = *token
    }

    return &Client{
        http: &http.Client{
            Timeout: timeout,
        },
        url: resolvedURL,
        authToken: resolvedToken,
    }, nil
}

func (c *Client) DoNothing() {
    fmt.Println("hello world")
}

//* Chat
//* Generate
//* Other

