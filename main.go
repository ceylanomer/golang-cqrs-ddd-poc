package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/infrastructure/persistence"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/interfaces/http/router"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger.Init()
	log := logger.GetLogger()
	defer log.Sync()

	// Initialize database connection
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	productRepo := persistence.NewPostgresRepository(db)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	})

	// Setup routes
	router.SetupProductRoutes(app, productRepo, productRepo)

	// Graceful shutdown channel
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatal("Error starting server", zap.Error(err))
		}
	}()

	log.Info("Server started on :8080")

	// Wait for shutdown signal
	<-shutdownChan
	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server gracefully stopped")
}
