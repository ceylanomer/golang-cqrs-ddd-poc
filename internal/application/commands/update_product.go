package commands

import (
	"context"
	"errors"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
)

type UpdateProductCommand struct {
	ID          uuid.UUID `json:"id" params:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	StockLevel  int       `json:"stock_level"`
	StockUnit   string    `json:"stock_unit"`
	Version     int       `json:"version"`
}

type UpdateProductHandler struct {
	repo product.Repository
}

func NewUpdateProductHandler(repo product.Repository) *UpdateProductHandler {
	return &UpdateProductHandler{repo: repo}
}

func (h *UpdateProductHandler) Handle(ctx context.Context, cmd *UpdateProductCommand) (*product.Product, error) {
	// Get existing product from repository
	existingProduct, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	// Check version for optimistic locking
	if existingProduct.Version() != cmd.Version {
		return nil, errors.New("product has been modified by another process")
	}

	// Update price if provided
	if cmd.Price > 0 {
		newPrice, err := product.NewPrice(cmd.Price, cmd.Currency)
		if err != nil {
			return nil, err
		}
		// Use domain logic to update price
		if err := existingProduct.UpdatePrice(newPrice); err != nil {
			return nil, err
		}
	}

	// Update stock if provided
	if cmd.StockLevel >= 0 {
		// Use domain logic to update stock
		if err := existingProduct.UpdateStock(cmd.StockLevel); err != nil {
			return nil, err
		}
	}

	// Persist updated product
	if err := h.repo.Update(ctx, existingProduct); err != nil {
		return nil, err
	}

	return existingProduct, nil
}
