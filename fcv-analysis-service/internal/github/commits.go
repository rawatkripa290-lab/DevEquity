package github

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type CommitItem struct {
	SHA     string `json:"sha"`
	HTMLURL string `json:"html_url"`
	Commit  struct {
		Message string `json:"message"`
		Author  struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
	Author *struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	} `json:"author"`
}

func (c *Client) ListCommits(ctx context.Context, owner, repo, branch string, limit int) ([]CommitItem, error) {
	if limit < 1 {
		return nil, fmt.Errorf("limit must be positive")
	}
	if limit > 1000 {
		limit = 1000
	}
	all := make([]CommitItem, 0, limit)
	for page := 1; len(all) < limit; page++ {
		perPage := 100
		if remaining := limit - len(all); remaining < perPage {
			perPage = remaining
		}
		endpoint := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/commits?per_page=%d&page=%d",
			url.PathEscape(owner), url.PathEscape(repo), perPage, page,
		)
		if branch != "" {
			endpoint += "&sha=" + url.QueryEscape(branch)
		}
		var batch []CommitItem
		if err := c.getJSON(ctx, endpoint, &batch); err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			break
		}
		all = append(all, batch...)
		if len(batch) < perPage {
			break
		}
	}
	return all, nil
}
