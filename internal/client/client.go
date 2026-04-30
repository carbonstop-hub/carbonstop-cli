// Package client provides an HTTP client for the Carbonstop Gateway API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/carbonstop/carbonstop-cli/internal/config"
)

const (
	userAgent    = "carbonstop-cli/0.2.0"
	maxRetries   = 3
	baseBackoff  = 1 * time.Second
)

// APIError represents an API-level error response.
type APIError struct {
	StatusCode int
	APICode    int    `json:"code"`
	Message    string `json:"msg"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}

// Client wraps the HTTP client for gateway API communication.
type Client struct {
	config     *config.Config
	httpClient *http.Client
}

// New creates a new Client from config.
func New(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// Request makes an HTTP request with retry logic and returns the raw status/body pair.
// Retries up to 3 times on 5xx and transport errors with exponential backoff (1s, 2s, 4s).
// 4xx errors are not retried.
func (c *Client) Request(method, path string, body io.Reader, query map[string]string, extraHeaders map[string]string) (int, string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	u, err := url.Parse(c.config.BaseURL + path)
	if err != nil {
		return 0, "", fmt.Errorf("invalid URL: %w", err)
	}

	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			if v != "" {
				q.Set(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}

	// Buffer body for replay on retry.
	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return 0, "", fmt.Errorf("read body: %w", err)
		}
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * baseBackoff
			time.Sleep(backoff)
		}

		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequest(strings.ToUpper(method), u.String(), reqBody)
		if err != nil {
			return 0, "", fmt.Errorf("request creation failed: %w", err)
		}

		req.Header.Set("X-API-Key", c.config.APIKey)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", userAgent)

		for k, v := range extraHeaders {
			req.Header.Set(k, v)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("transport error: %w", err)
			continue
		}

		respBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if readErr != nil {
			lastErr = fmt.Errorf("read error: %w", readErr)
			continue
		}

		// 429 (rate limit) and 5xx are retryable — check before 4xx.
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			lastErr = &APIError{StatusCode: resp.StatusCode, Message: "retryable server error"}
			continue
		}

		// Other 4xx are client errors — do not retry.
		if resp.StatusCode >= 400 {
			return resp.StatusCode, string(respBytes), nil
		}

		return resp.StatusCode, string(respBytes), nil
	}

	// All retries exhausted.
	return 0, "", fmt.Errorf("request failed after %d retries: %w", maxRetries+1, lastErr)
}

// Get makes a GET request.
func (c *Client) Get(path string, query map[string]string) (int, string, error) {
	return c.Request("GET", path, nil, query, nil)
}

// Post makes a POST request with a JSON body.
func (c *Client) Post(path string, jsonBody []byte) (int, string, error) {
	var body io.Reader
	var extraHeaders map[string]string
	if jsonBody != nil {
		body = strings.NewReader(string(jsonBody))
		extraHeaders = map[string]string{"Content-Type": "application/json;charset=utf-8"}
	}
	return c.Request("POST", path, body, nil, extraHeaders)
}

// Params holds common query parameters for API calls.
type Params struct {
	PageNum    int
	PageSize   int
	Status     int
	Search     string
	ProductID  int
	AccountID  int
	GroupType  int
	Lang       string
	AccountStatus int
}

// ParseResponse parses a JSON response body and checks for errors.
func ParseResponse(status int, text string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		// Return raw text as msg
		return map[string]interface{}{"code": status, "msg": text}, nil
	}
	if status < 200 || status >= 300 {
		msg := "unknown error"
		if m, ok := data["msg"]; ok {
			msg = fmt.Sprint(m)
		}
		return data, &APIError{StatusCode: status, Message: msg}
	}
	return data, nil
}
