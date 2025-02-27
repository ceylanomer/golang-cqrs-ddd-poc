package commands

import (
	"context"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
)

type DeleteProductCommand struct {
	ID uuid.UUID
}

type DeleteProductHandler struct {
	repo product.Repository
}

func NewDeleteProductHandler(repo product.Repository) *DeleteProductHandler {
	return &DeleteProductHandler{
		repo: repo,
	}
}

func (h *DeleteProductHandler) Handle(ctx context.Context, cmd *DeleteProductCommand) (*product.Product, error) {
	product, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return nil, err
	}

	return product, nil
}
