package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/readspark/backend/internal/config"
	"github.com/readspark/backend/internal/database"
	"github.com/readspark/backend/internal/handler"
	"github.com/readspark/backend/internal/middleware"
	"github.com/readspark/backend/internal/repository"
	"github.com/readspark/backend/internal/scheduler"
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
	progressRepo := repository.NewProgressRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWT)
	searcher := service.NewPGFullTextSearch(articleRepo)
	articleService := service.NewArticleService(articleRepo, searcher)
	progressService := service.NewProgressService(progressRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	dictionaryService := service.NewDictionaryService()

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	articleHandler := handler.NewArticleHandler(articleService)
	progressHandler := handler.NewProgressHandler(progressService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	dictionaryHandler := handler.NewDictionaryHandler(dictionaryService)

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
		api.GET("/dictionary/:word", dictionaryHandler.Lookup)

		// Protected
		authenticated := api.Group("/")
		authenticated.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			authenticated.GET("/articles/:id", articleHandler.GetArticle)
			authenticated.POST("/progress", progressHandler.Sync)
			authenticated.GET("/progress", progressHandler.List)
			authenticated.POST("/subscriptions", subscriptionHandler.Create)
			authenticated.GET("/subscriptions/status", subscriptionHandler.Status)
		}
	}

	jobScheduler, err := scheduler.New()
	if err != nil {
		slog.Error("failed to initialize scheduler", "error", err)
		os.Exit(1)
	}
	jobScheduler.Start()
	defer jobScheduler.Stop()

	slog.Info("server starting", "port", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
