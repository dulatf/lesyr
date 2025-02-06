package server

import (
	"fmt"

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

func NewServer(cfg *config.Config) (*fiber.App, error) {
	// Initialize database connection
	dbPool, err := database.NewDBPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Initialize repositories
	// userRepo := postgres.NewUserRepository(dbPool)
	feedRepo := postgres.NewFeedRepository(dbPool)
	articleRepo := postgres.NewArticleRepository(dbPool)

	// Initialize handlers
	feedFetcher := service.NewFeedFetcher(feedRepo, articleRepo)

	// Handlers
	feedHandler := handler.NewFeedHandler(feedRepo, feedFetcher)
	articleHandler := handler.NewArticleHandler(articleRepo, feedRepo)

	app := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check
	app.Get("/health", handler.HealthCheck)

	// API routes will be grouped under /api/v1
	v1 := app.Group("/api/v1")

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// Feed routes (protected)
	feeds := v1.Group("/feeds")
	feeds.Use(handler.AuthMiddleware)
	feeds.Get("/", feedHandler.ListFeeds)
	feeds.Post("/", feedHandler.CreateFeed)
	feeds.Get("/:id", feedHandler.GetFeed)
	feeds.Delete("/:id", feedHandler.DeleteFeed)
	feeds.Post("/:id/refresh", feedHandler.RefreshFeed)

	// Article routes (protected)
	articles := v1.Group("/articles")
	articles.Use(handler.AuthMiddleware)
	articles.Get("/unread", articleHandler.GetUnreadArticles)
	articles.Get("/:id", articleHandler.GetArticle)
	articles.Post("/:id/read", articleHandler.MarkArticleRead)

	return app, nil
}
