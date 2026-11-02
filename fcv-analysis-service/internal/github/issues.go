package github

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type Issue struct {
	Number  int64  `json:"number"`
	Title   string `json:"title"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at"`

	Comments int `json:"comments"`

	User struct {
		Login string `json:"login"`
	} `json:"user"`

	// GitHub's issues endpoint also returns pull requests.
	PullRequest *struct {
		URL string `json:"url"`
	} `json:"pull_request"`
}

func (c *Client) ListIssues(
	ctx context.Context,
	owner, repo string,
	limit int,
) ([]Issue, error) {

	if limit < 1 {
		return nil, fmt.Errorf("limit must be positive")
	}

	if limit > 1000 {
		limit = 1000
	}

	result := make([]Issue, 0, limit)

	// GitHub's issues endpoint includes PRs.
	// We filter those out here.
	for page := 1; len(result) < limit; page++ {

		perPage := 100

		if remaining := limit - len(result); remaining < perPage {
			perPage = remaining
		}

		endpoint := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/issues?state=all&per_page=%d&page=%d",
			url.PathEscape(owner),
			url.PathEscape(repo),
			perPage,
			page,
		)

		var batch []Issue

		if err := c.getJSON(ctx, endpoint, &batch); err != nil {
			return nil, fmt.Errorf("fetch issues: %w", err)
		}

		if len(batch) == 0 {
			break
		}

		for _, issue := range batch {
			if issue.PullRequest != nil {
				continue
			}

			result = append(result, issue)

			if len(result) >= limit {
				break
			}
		}

		if len(batch) < perPage {
			break
		}
	}

	return result, nil
}
