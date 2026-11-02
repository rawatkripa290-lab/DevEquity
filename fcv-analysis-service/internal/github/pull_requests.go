package github

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type PullRequest struct {
	Number  int64  `json:"number"`
	Title   string `json:"title"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	MergedAt  *time.Time `json:"merged_at"`

	User struct {
		Login string `json:"login"`
	} `json:"user"`

	Additions    int `json:"additions"`
	Deletions    int `json:"deletions"`
	ChangedFiles int `json:"changed_files"`
	Comments     int `json:"comments"`
}

func (c *Client) ListPullRequests(
	ctx context.Context,
	owner, repo string,
	limit int,
) ([]PullRequest, error) {

	if limit < 1 {
		return nil, fmt.Errorf("limit must be positive")
	}

	if limit > 1000 {
		limit = 1000
	}

	result := make([]PullRequest, 0, limit)

	for page := 1; len(result) < limit; page++ {

		perPage := 100

		if remaining := limit - len(result); remaining < perPage {
			perPage = remaining
		}

		endpoint := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/pulls?state=all&per_page=%d&page=%d",
			url.PathEscape(owner),
			url.PathEscape(repo),
			perPage,
			page,
		)

		var batch []PullRequest

		if err := c.getJSON(ctx, endpoint, &batch); err != nil {
			return nil, fmt.Errorf("fetch pull requests: %w", err)
		}

		if len(batch) == 0 {
			break
		}

		result = append(result, batch...)

		if len(batch) < perPage {
			break
		}
	}

	return result, nil
}

// GitHub's list endpoint does not include additions,
// deletions and changed_files in every response.
// Fetch individual PR details when those values are needed.
func (c *Client) GetPullRequest(
	ctx context.Context,
	owner, repo string,
	number int64,
) (PullRequest, error) {

	var pr PullRequest

	endpoint := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/pulls/%d",
		url.PathEscape(owner),
		url.PathEscape(repo),
		number,
	)

	if err := c.getJSON(ctx, endpoint, &pr); err != nil {
		return pr, fmt.Errorf("fetch pull request %d: %w", number, err)
	}

	return pr, nil
}
