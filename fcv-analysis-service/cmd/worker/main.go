package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/yourusername/fair-contribution-view/backend/internal/ai"
	"github.com/yourusername/fair-contribution-view/backend/internal/config"
	"github.com/yourusername/fair-contribution-view/backend/internal/db"
	"github.com/yourusername/fair-contribution-view/backend/internal/github"
	"github.com/yourusername/fair-contribution-view/backend/internal/queue"
	"github.com/yourusername/fair-contribution-view/backend/internal/worker"
)

func main() {
	cfg := config.Load()
	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	mux := asynq.NewServeMux()
	handler := &worker.IngestWorker{
		DB:          pool,
		GitHub:      github.NewClient(cfg.GitHubToken),
		AI:          ai.NewClient(cfg.AIServiceURL),
		CommitLimit: cfg.CommitLimit,
	}
	mux.HandleFunc(queue.TypeIngestRepository, handler.Handle)

	srv := queue.NewServer(cfg.RedisAddr, cfg.RedisPass)
	log.Println("Go ingestion worker started")
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
