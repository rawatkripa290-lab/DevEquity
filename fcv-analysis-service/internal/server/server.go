package server

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/fair-contribution-view/backend/internal/handlers"
)

func NewRouter(h *handlers.AnalysisHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	v1 := r.Group("/api/v1")
	v1.POST("/analyses", h.Create)
	v1.GET("/analyses/:id", h.Get)
	v1.GET("/analyses/:id/report", h.Report)
	return r
}
