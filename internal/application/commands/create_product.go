package commands

import (
	"context"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/gofiber/fiber/v2"
)

type CreateProductCommand struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	StockLevel  int     `json:"stock_level"`
	StockUnit   string  `json:"stock_unit"`
}

type CreateProductHandler struct {
	repo product.Repository
}

func NewCreateProductHandler(repo product.Repository) *CreateProductHandler {
	return &CreateProductHandler{repo: repo}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd *CreateProductCommand) (*product.Product, error) {
	// Create value objects using domain logic
	price, err := product.NewPrice(cmd.Price, cmd.Currency)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	stock, err := product.NewStock(cmd.StockLevel, cmd.StockUnit)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create new product using domain factory
	newProduct, err := product.NewProduct(cmd.Name, cmd.Description, price, stock)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Persist using repository
	if err := h.repo.Save(ctx, newProduct); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return newProduct, nil
}
