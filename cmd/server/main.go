package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dulatf/lesyr/internal/config"
	"github.com/dulatf/lesyr/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Create channel for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		log.Printf("Starting server on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-c

	// Cleanup resources
	log.Println("Shutting down server...")
	server.Cleanup()

	if err := app.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
