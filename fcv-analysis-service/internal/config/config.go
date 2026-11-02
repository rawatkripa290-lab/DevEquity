package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	RedisAddr    string
	RedisPass    string
	GitHubToken  string
	AIServiceURL string
	CommitLimit  int
}

func Load() Config {
	limit, err := strconv.Atoi(getenv("COMMIT_LIMIT", "300"))
	if err != nil || limit < 1 {
		limit = 300
	}
	return Config{
		Port:         getenv("PORT", "8080"),
		DatabaseURL:  getenv("DATABASE_URL", "postgres://fairview:fairview_dev_password@localhost:5432/fairview?sslmode=disable"),
		RedisAddr:    getenv("REDIS_ADDR", "localhost:6379"),
		RedisPass:    os.Getenv("REDIS_PASSWORD"),
		GitHubToken:  os.Getenv("GITHUB_TOKEN"),
		AIServiceURL: getenv("AI_SERVICE_URL", "http://ai-service:8000"),
		CommitLimit:  limit,
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
