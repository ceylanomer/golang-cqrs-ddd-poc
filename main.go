package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/application/client"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/infrastructure/persistence"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/interfaces/http/router"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/config"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/logger"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/tracer"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	recover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func main() {
	// Initialize logger
	logger.Init()
	log := logger.GetLogger()
	defer log.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		zap.L().Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize clients
	transport := client.NewTransport()
	noRetryClient := client.NewHttpClient(transport)
	retryableClient := client.NewRetryableClient(transport)

	// Initialize tracer
	tp := tracer.InitTracer(cfg.Jaeger)

	// Create custom GORM logger with Zap
	gormLogger := persistence.NewGormZapLogger(log)

	// Initialize database connection with GORM
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger:      gormLogger,
		QueryFields: true, // Enable query fields for better tracing
	})
	if err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}

	// Register OpenTelemetry callbacks
	err = db.Use(tracing.NewPlugin(tracing.WithTracerProvider(tp)))
	if err != nil {
		zap.L().Fatal("Failed to register GORM tracing", zap.Error(err))
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&persistence.ProductModel{}); err != nil {
		zap.L().Fatal("Failed to migrate database schema", zap.Error(err))
	}

	// Initialize repositories
	productRepo := persistence.NewProductRepository(db)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	})

	// Middleware order is important for tracing
	app.Use(recover.New())
	// Add OpenTelemetry middleware first to create the parent span
	app.Use(otelfiber.Middleware(otelfiber.WithNext(func(c *fiber.Ctx) bool {
		// Skip tracing for metrics endpoint
		return c.Path() == "/metrics"
	})))
	// Then add logging middleware
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: log,
		SkipURIs: []string{"/metrics"},
	}))

	// Setup routes
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	router.SetupProductRoutes(app, productRepo, productRepo, noRetryClient, retryableClient)

	// Graceful shutdown channel
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		if err := app.Listen(serverAddr); err != nil {
			zap.L().Fatal("Error starting server", zap.Error(err))
		}
	}()

	zap.L().Info("Server started", zap.Int("port", cfg.Server.Port))

	// Wait for shutdown signal
	<-shutdownChan
	zap.L().Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Close database connection
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("Error getting underlying *sql.DB", zap.Error(err))
	} else {
		if err := sqlDB.Close(); err != nil {
			zap.L().Error("Error closing database connection", zap.Error(err))
		}
	}

	zap.L().Info("Server gracefully stopped")
}
