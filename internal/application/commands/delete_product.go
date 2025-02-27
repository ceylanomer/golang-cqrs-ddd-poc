package commands

import (
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

func (h *DeleteProductHandler) Handle(cmd DeleteProductCommand) error {
	return h.repo.Delete(cmd.ID)
}
