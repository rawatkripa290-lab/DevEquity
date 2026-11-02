package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/fair-contribution-view/backend/internal/config"
	"github.com/yourusername/fair-contribution-view/backend/internal/db"
	"github.com/yourusername/fair-contribution-view/backend/internal/github"
	"github.com/yourusername/fair-contribution-view/backend/internal/handlers"
	"github.com/yourusername/fair-contribution-view/backend/internal/queue"
	"github.com/yourusername/fair-contribution-view/backend/internal/server"
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
