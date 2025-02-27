package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/application/commands"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/application/queries"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/infrastructure/persistence"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/interfaces/http/handlers"
	applogger "github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/logger"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	applogger.Init()
	log := applogger.GetLogger()
	defer zap.L().Sync()
	defer log.Sync()

	// Initialize repository
	productRepo := persistence.NewMemoryProductRepository()

	// Add some sample data
	sampleProduct := &product.Product{
		ID:          uuid.New(),
		Name:        "Sample Product",
		Description: "This is a sample product",
		Price:       99.99,
	}
	productRepo.Save(sampleProduct)
	//log.Info("Created sample product", zap.String("id", sampleProduct.ID.String()))
	zap.L().Info("Created sample product", zap.String("id", sampleProduct.ID.String()))

	// Initialize handlers
	getProductHandler := queries.NewGetProductHandler(productRepo)
	createProductHandler := commands.NewCreateProductHandler(productRepo)
	updateProductHandler := commands.NewUpdateProductHandler(productRepo)
	deleteProductHandler := commands.NewDeleteProductHandler(productRepo)

	// Initialize HTTP handlers
	productHandler := handlers.NewProductHandler(
		getProductHandler,
		createProductHandler,
		updateProductHandler,
		deleteProductHandler,
	)

	// Setup Fiber with custom config
	app := fiber.New(fiber.Config{
		// Add graceful shutdown timeout
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	})

	// Use fiberzap to log requests
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: log,
	}))

	// Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	products := v1.Group("/products")
	products.Get("/:id", productHandler.GetProduct)
	products.Post("/", productHandler.CreateProduct)
	products.Put("/:id", productHandler.UpdateProduct)
	products.Delete("/:id", productHandler.DeleteProduct)

	// Create channel for shutdown signals
	shutdownChan := make(chan os.Signal, 1)
	// Catch SIGINT and SIGTERM signals
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(":8080"); err != nil {
			zap.L().Fatal("Error starting server", zap.Error(err))
		}
	}()

	zap.L().Info("Server started on :8080")

	// Wait for shutdown signal
	<-shutdownChan
	zap.L().Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := app.ShutdownWithContext(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown", zap.Error(err))
	}

	zap.L().Info("Server gracefully stopped")
}
