package github

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type Review struct {
	ID    int64  `json:"id"`
	State string `json:"state"`
	Body  string `json:"body"`

	HTMLURL string `json:"html_url"`

	SubmittedAt *time.Time `json:"submitted_at"`

	User struct {
		Login string `json:"login"`
	} `json:"user"`
}

func (c *Client) ListReviews(
	ctx context.Context,
	owner, repo string,
	prNumber int64,
) ([]Review, error) {

	result := make([]Review, 0)

	for page := 1; page <= 10; page++ {

		endpoint := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/pulls/%d/reviews?per_page=100&page=%d",
			url.PathEscape(owner),
			url.PathEscape(repo),
			prNumber,
			page,
		)

		var batch []Review

		if err := c.getJSON(ctx, endpoint, &batch); err != nil {
			return nil, fmt.Errorf(
				"fetch reviews for PR %d: %w",
				prNumber,
				err,
			)
		}

		if len(batch) == 0 {
			break
		}

		result = append(result, batch...)

		if len(batch) < 100 {
			break
		}
	}

	return result, nil
}
