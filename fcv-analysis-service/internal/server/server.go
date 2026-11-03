
package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rawatkripa290-lab/DevEquity/fcv-analysis-service/internal/handlers"
)

func NewRouter(h *handlers.AnalysisHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://dev-equity.vercel.app",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	v1.POST("/analyses", h.Create)
	v1.GET("/analyses/:id", h.Get)
	v1.GET("/analyses/:id/report", h.Report)

	return r
}
