package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/idcim/aihuman/backend/internal/api"
	"github.com/idcim/aihuman/backend/internal/config"
	"github.com/idcim/aihuman/backend/internal/video"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	store := video.NewMemoryStore()
	service := video.NewService(store, logger)
	handler := video.NewHandler(service)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(logger))

	router.GET("/healthz", func(c *gin.Context) {
		api.Success(c, gin.H{"status": "ok", "service": "backend"})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/video/create", handler.Create)
		v1.GET("/video/:id", handler.Get)
	}

	addr := ":" + cfg.Port
	logger.Info("backend listening", "addr", addr)
	if err := router.Run(addr); err != nil && err != http.ErrServerClosed {
		logger.Error("backend stopped", "error", err)
		os.Exit(1)
	}
}

func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logger.Info(
			"http_request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"request_id", c.GetHeader("X-Request-Id"),
		)
	}
}
