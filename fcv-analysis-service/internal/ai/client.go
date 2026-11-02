package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Commit struct {
	SHA         string   `json:"sha"`
	Message     string   `json:"message"`
	Files       []string `json:"files"`
	DiffSummary string   `json:"diff_summary"`
}
type Result struct {
	SHA         string  `json:"sha"`
	Category    string  `json:"category"`
	Confidence  float64 `json:"confidence"`
	Explanation string  `json:"explanation"`
}
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL, HTTP: &http.Client{Timeout: 12 * time.Second}}
}
func (c *Client) Classify(ctx context.Context, commits []Commit) ([]Result, error) {
	body, err := json.Marshal(map[string]any{"commits": commits})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/internal/classify", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("AI service returned HTTP %d", resp.StatusCode)
	}
	var out []Result
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
