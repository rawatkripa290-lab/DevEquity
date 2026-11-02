package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Classify(ctx context.Context, commits []Commit) ([]Result, error) {
	if len(commits) == 0 {
		return []Result{}, nil
	}

	// Python expects files to be a list, never null.
	for i := range commits {
		if commits[i].Files == nil {
			commits[i].Files = []string{}
		}
	}

	body, err := json.Marshal(map[string]any{
		"commits": commits,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal AI request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseURL+"/internal/classify",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create AI request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call AI service: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read AI response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"AI service returned HTTP %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var results []Result
	if err := json.Unmarshal(responseBody, &results); err != nil {
		return nil, fmt.Errorf("decode AI response: %w", err)
	}

	if len(results) != len(commits) {
		return nil, fmt.Errorf(
			"AI returned %d classifications for %d commits",
			len(results),
			len(commits),
		)
	}

	return results, nil
}
