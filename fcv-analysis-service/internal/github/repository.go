package github

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var repoPart = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

type Repository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	DefaultBranch string `json:"default_branch"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func ParseRepositoryURL(raw string) (owner, name string, err error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", "", fmt.Errorf("invalid repository URL")
	}
	if u.Scheme != "https" || !strings.EqualFold(u.Hostname(), "github.com") ||
		u.User != nil || u.Port() != "" {
		return "", "", fmt.Errorf("only HTTPS github.com repository URLs are accepted")
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected URL format https://github.com/owner/repository")
	}
	owner = parts[0]
	name = strings.TrimSuffix(parts[1], ".git")
	if !repoPart.MatchString(owner) || !repoPart.MatchString(name) {
		return "", "", fmt.Errorf("invalid owner or repository name")
	}
	return owner, name, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, name string) (Repository, error) {
	var repo Repository
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s", url.PathEscape(owner), url.PathEscape(name))
	err := c.getJSON(ctx, endpoint, &repo)
	return repo, err
}
