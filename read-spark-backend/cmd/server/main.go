package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/readspark/backend/internal/config"
	"github.com/readspark/backend/internal/database"
	"github.com/readspark/backend/internal/domain"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Article{},
		&domain.Subscription{},
		&domain.ReadingProgress{},
		&domain.Annotation{},
	); err != nil {
		slog.Error("failed to migrate", "error", err)
		os.Exit(1)
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	slog.Info("server starting", "port", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
