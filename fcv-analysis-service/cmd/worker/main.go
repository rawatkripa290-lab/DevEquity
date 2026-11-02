package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/ai"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/config"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/db"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/github"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/queue"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/worker"
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
