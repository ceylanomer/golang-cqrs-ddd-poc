package router

import (
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/application/client"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/application/commands"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/application/queries"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(
	app *fiber.App,
	writeRepo product.Repository,
	readRepo product.ReadOnlyRepository,
	noRetryClient client.CustomHttpClient,
	retryableClient client.CustomRetryableClient,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	products := v1.Group("/products")

	// Command handlers
	createHandler := commands.NewCreateProductHandler(writeRepo)
	updateHandler := commands.NewUpdateProductHandler(writeRepo)
	statusHandler := commands.NewChangeProductStatusHandler(writeRepo)
	deleteHandler := commands.NewDeleteProductHandler(writeRepo)
	// Query handlers
	getHandler := queries.NewGetProductHandler(readRepo)
	listHandler := queries.NewListProductsHandler(readRepo)

	// Routes
	products.Post("/", handler.Handler(createHandler))
	products.Get("/:id", handler.Handler(getHandler))
	products.Put("/:id", handler.Handler(updateHandler))
	products.Put("/:id/status", handler.Handler(statusHandler))
	products.Get("/", handler.Handler(listHandler))
	products.Delete("/:id", handler.Handler(deleteHandler))
}
