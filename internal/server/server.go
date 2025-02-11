package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dulatf/lesyr/internal/config"
	"github.com/dulatf/lesyr/internal/database"
	"github.com/dulatf/lesyr/internal/handler"
	"github.com/dulatf/lesyr/internal/repository/postgres"
	"github.com/dulatf/lesyr/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Global scheduler for cleanup
var (
	scheduler     *service.FeedScheduler
	shutdownOnce  sync.Once
	shutdownHooks []func() error
)

// AddShutdownHook adds a function to be called on server shutdown
func AddShutdownHook(fn func() error) {
	shutdownHooks = append(shutdownHooks, fn)
}

func Cleanup() {
	shutdownOnce.Do(func() {
		for _, hook := range shutdownHooks {
			if err := hook(); err != nil {
				fmt.Printf("Error during shutdown: %v\n", err)
			}
		}
	})
}

func NewServer(cfg *config.Config) (*fiber.App, error) {
	// Initialize database connection
	dbPool, err := database.NewDBPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Load OAuth configuration
	oauthConfig, err := config.LoadOAuthConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load OAuth config: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(dbPool)
	feedRepo := postgres.NewFeedRepository(dbPool)
	articleRepo := postgres.NewArticleRepository(dbPool)

	// Initialize handlers
	feedFetcher := service.NewFeedFetcher(feedRepo, articleRepo)
	feedScheduler := service.NewFeedScheduler(
		feedRepo,
		feedFetcher,
		10,             // number of workers
		15*time.Minute, // refresh interval
	)

	// Start the scheduler
	if err := feedScheduler.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to start feed scheduler: %v", err)
	}

	// Add cleanup hook for scheduler
	AddShutdownHook(func() error {
		if scheduler != nil {
			scheduler.Stop()
		}
		return nil
	})

	// Handlers
	authHandler := handler.NewAuthHandler(userRepo, oauthConfig, cfg.JWTSecret)
	feedHandler := handler.NewFeedHandler(feedRepo, feedFetcher)
	articleHandler := handler.NewArticleHandler(articleRepo, feedRepo)

	app := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", handler.HealthCheck)

	// API routes will be grouped under /api/v1
	v1 := app.Group("/api/v1")

	// Auth routes
	auth := v1.Group("/auth")
	auth.Get("/github", authHandler.InitiateGitHubAuth)
	auth.Get("/callback", authHandler.HandleGitHubCallback)

	api := v1.Group("")
	api.Use(authHandler.AuthMiddleware)

	auth_me := auth.Group("")
	auth_me.Use(authHandler.AuthMiddleware)
	auth_me.Get("/me", authHandler.GetAuthenticatedUser)

	// Feed routes (protected)
	feeds := api.Group("/feeds")
	feeds.Get("/", feedHandler.ListFeeds)
	feeds.Post("/", feedHandler.CreateFeed)
	feeds.Get("/:id", feedHandler.GetFeed)
	feeds.Delete("/:id", feedHandler.DeleteFeed)
	feeds.Post("/:id/refresh", feedHandler.RefreshFeed)
	feeds.Get("/:id/articles", articleHandler.GetFeedArticles)

	// Article routes (protected)
	articles := api.Group("/articles")
	articles.Get("/unread", articleHandler.GetUnreadArticles)
	articles.Get("/:id", articleHandler.GetArticle)
	articles.Post("/:id/read", articleHandler.MarkArticleRead)

	return app, nil
}
