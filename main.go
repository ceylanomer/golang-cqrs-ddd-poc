package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/infrastructure/persistence"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/interfaces/http/router"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/config"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize logger
	logger.Init()
	log := logger.GetLogger()
	defer log.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize database connection with GORM
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&persistence.ProductModel{}); err != nil {
		log.Fatal("Failed to migrate database schema", zap.Error(err))
	}

	// Initialize repositories
	productRepo := persistence.NewProductRepository(db)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	})

	// Setup routes
	router.SetupProductRoutes(app, productRepo, productRepo)

	// Graceful shutdown channel
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		if err := app.Listen(serverAddr); err != nil {
			log.Fatal("Error starting server", zap.Error(err))
		}
	}()

	log.Info("Server started", zap.Int("port", cfg.Server.Port))

	// Wait for shutdown signal
	<-shutdownChan
	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Close database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("Error getting underlying *sql.DB", zap.Error(err))
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Error("Error closing database connection", zap.Error(err))
		}
	}

	log.Info("Server gracefully stopped")
}
