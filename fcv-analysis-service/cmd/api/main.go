
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/ai"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/config"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/db"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/github"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/handlers"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/queue"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/server"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/worker"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// Connect to PostgreSQL
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Connect to Redis for adding jobs
	q := queue.NewClient(cfg.RedisAddr, cfg.RedisPass)
	defer q.Close()

	// Configure the API handler
	h := &handlers.AnalysisHandler{
		DB:     pool,
		Queue:  q,
		GitHub: github.NewClient(cfg.GitHubToken),
	}

	// Configure the background worker
	mux := asynq.NewServeMux()
	handler := &worker.IngestWorker{
		DB:          pool,
		GitHub:      github.NewClient(cfg.GitHubToken),
		AI:          ai.NewClient(cfg.AIServiceURL),
		CommitLimit: cfg.CommitLimit,
	}
	mux.HandleFunc(queue.TypeIngestRepository, handler.Handle)

	workerServer := queue.NewServer(cfg.RedisAddr, cfg.RedisPass)

	// Start worker in the background
	go func() {
		log.Println("Background worker started")
		if err := workerServer.Run(mux); err != nil {
			log.Printf("WORKER ERROR: %v", err)
		}
	}()

	// Start HTTP API
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.NewRouter(h),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("API listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("API ERROR: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down API and worker...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("API shutdown error: %v", err)
	}

	workerServer.Shutdown()
}
