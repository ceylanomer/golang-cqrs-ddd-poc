package commands

import (
	"fmt"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
)

type UpdateProductCommand struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}

type UpdateProductHandler struct {
	repo product.Repository
}

func NewUpdateProductHandler(repo product.Repository) *UpdateProductHandler {
	return &UpdateProductHandler{
		repo: repo,
	}
}

func (h *UpdateProductHandler) Handle(cmd UpdateProductCommand) (*product.Product, error) {
	// First check if product exists
	existing, err := h.repo.GetByID(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Update the product
	existing.Name = cmd.Name
	existing.Description = cmd.Description
	existing.Price = cmd.Price

	if err := h.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}
