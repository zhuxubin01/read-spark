package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/readspark/backend/internal/config"
	"github.com/readspark/backend/internal/database"
	"github.com/readspark/backend/internal/domain"
	"github.com/readspark/backend/internal/handler"
	"github.com/readspark/backend/internal/middleware"
	"github.com/readspark/backend/internal/repository"
	"github.com/readspark/backend/internal/service"
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

	// Repositories
	userRepo := repository.NewUserRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWT)
	searcher := service.NewPGFullTextSearch(articleRepo)
	articleService := service.NewArticleService(articleRepo, searcher)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	articleHandler := handler.NewArticleHandler(articleService)

	// Routes
	api := r.Group("/api/v1")
	{
		// Public
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}

		api.GET("/articles/daily", articleHandler.GetDaily)
		api.GET("/articles", articleHandler.ListArticles)

		// Protected
		authenticated := api.Group("/")
		authenticated.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			authenticated.GET("/articles/:id", articleHandler.GetArticle)
		}
	}

	slog.Info("server starting", "port", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
