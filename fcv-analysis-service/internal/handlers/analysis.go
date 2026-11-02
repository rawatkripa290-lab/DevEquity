package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/github"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/queue"
)

type AnalysisHandler struct {
	DB     *pgxpool.Pool
	Queue  *asynq.Client
	GitHub *github.Client
}

type createAnalysisRequest struct {
	RepositoryURL string `json:"repository_url" binding:"required"`
}

func (h *AnalysisHandler) Create(c *gin.Context) {
	var req createAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repository_url is required"})
		return
	}
	owner, name, err := github.ParseRepositoryURL(req.RepositoryURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()
	repo, err := h.GitHub.GetRepository(ctx, owner, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not retrieve repository from GitHub"})
		return
	}

	var repositoryID, analysisID string
	err = h.DB.QueryRow(ctx, `
		INSERT INTO repositories (github_repo_id, owner, name, full_name, default_branch)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (github_repo_id) DO UPDATE SET
			owner=EXCLUDED.owner, name=EXCLUDED.name,
			full_name=EXCLUDED.full_name, default_branch=EXCLUDED.default_branch
		RETURNING id
	`, repo.ID, repo.Owner.Login, repo.Name, repo.FullName, repo.DefaultBranch).Scan(&repositoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save repository"})
		return
	}

	err = h.DB.QueryRow(ctx, `
		INSERT INTO analyses (repository_id, status)
		VALUES ($1, 'queued') RETURNING id
	`, repositoryID).Scan(&analysisID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create analysis"})
		return
	}

	payload, err := queue.EncodeIngest(queue.IngestPayload{
		AnalysisID: analysisID, Owner: repo.Owner.Login, Repo: repo.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not encode job"})
		return
	}
	_, err = h.Queue.Enqueue(asynq.NewTask(queue.TypeIngestRepository, payload),
		asynq.MaxRetry(3), asynq.Timeout(10*time.Minute))
	if err != nil {
		_, _ = h.DB.Exec(context.Background(),
			`UPDATE analyses SET status='failed', error_message='queue unavailable' WHERE id=$1`,
			analysisID)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "analysis queue unavailable"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"analysis_id": analysisID,
		"status":      "queued",
	})
}

func (h *AnalysisHandler) Get(c *gin.Context) {
	var status, progress string
	var errMsg *string
	err := h.DB.QueryRow(c.Request.Context(),
		`SELECT status, progress::text, error_message FROM analyses WHERE id=$1`,
		c.Param("id")).Scan(&status, &progress, &errMsg)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "analysis not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"analysis_id": c.Param("id"),
		"status":      status,
		"progress":    progress,
		"error":       errMsg,
	})
}

func (h *AnalysisHandler) Report(c *gin.Context) {
	var raw []byte
	err := h.DB.QueryRow(c.Request.Context(), `SELECT report FROM analysis_reports WHERE analysis_id=$1`, c.Param("id")).Scan(&raw)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found; analysis may still be processing"})
		return
	}
	var report json.RawMessage = raw
	c.Data(http.StatusOK, "application/json; charset=utf-8", report)
}
