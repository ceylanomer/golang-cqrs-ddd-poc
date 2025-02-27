package handlers

import (
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/application/commands"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/application/queries"
	"github.com/ceylanomer/golang-cqrs-ddd-poc/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductHandler struct {
	getProductHandler    *queries.GetProductHandler
	createProductHandler *commands.CreateProductHandler
	updateProductHandler *commands.UpdateProductHandler
	deleteProductHandler *commands.DeleteProductHandler
}

func NewProductHandler(
	getProductHandler *queries.GetProductHandler,
	createProductHandler *commands.CreateProductHandler,
	updateProductHandler *commands.UpdateProductHandler,
	deleteProductHandler *commands.DeleteProductHandler,
) *ProductHandler {
	return &ProductHandler{
		getProductHandler:    getProductHandler,
		createProductHandler: createProductHandler,
		updateProductHandler: updateProductHandler,
		deleteProductHandler: deleteProductHandler,
	}
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.Error("Invalid product ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	query := queries.GetProductQuery{ID: id}
	product, err := h.getProductHandler.Handle(query)
	if err != nil {
		logger.Error("Failed to get product", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var cmd commands.CreateProductCommand
	if err := c.BodyParser(&cmd); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	product, err := h.createProductHandler.Handle(cmd)
	if err != nil {
		logger.Error("Failed to create product", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.Error("Invalid product ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var cmd commands.UpdateProductCommand
	if err := c.BodyParser(&cmd); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	cmd.ID = id

	product, err := h.updateProductHandler.Handle(cmd)
	if err != nil {
		logger.Error("Failed to update product", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.JSON(product)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.Error("Invalid product ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	cmd := commands.DeleteProductCommand{ID: id}
	if err := h.deleteProductHandler.Handle(cmd); err != nil {
		logger.Error("Failed to delete product", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
