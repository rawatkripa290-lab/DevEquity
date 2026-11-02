package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/config"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/db"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/github"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/handlers"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/queue"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/server"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	q := queue.NewClient(cfg.RedisAddr, cfg.RedisPass)
	defer q.Close()

	h := &handlers.AnalysisHandler{
		DB:     pool,
		Queue:  q,
		GitHub: github.NewClient(cfg.GitHubToken),
	}
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.NewRouter(h),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("API listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
