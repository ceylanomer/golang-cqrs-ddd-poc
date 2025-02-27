package commands

import (
	"context"
	"errors"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
)

type ChangeProductStatusCommand struct {
	ID      uuid.UUID `json:"id" params:"id"`
	Action  string    `json:"action"` // "activate", "deactivate", "discontinue"
	Version int       `json:"version"`
}

type ChangeProductStatusResponse struct {
	ID     uuid.UUID             `json:"id"`
	Status product.ProductStatus `json:"status"`
}

type ChangeProductStatusHandler struct {
	repo product.Repository
}

func NewChangeProductStatusHandler(repo product.Repository) *ChangeProductStatusHandler {
	return &ChangeProductStatusHandler{repo: repo}
}

func (h *ChangeProductStatusHandler) Handle(ctx context.Context, cmd *ChangeProductStatusCommand) (*ChangeProductStatusResponse, error) {
	existingProduct, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	
	if existingProduct.Version() != cmd.Version {
		return nil, errors.New("product has been modified by another process")
	}

	var actionErr error
	switch cmd.Action {
	case "activate":
		actionErr = existingProduct.Activate()
	case "deactivate":
		actionErr = existingProduct.Deactivate()
	case "discontinue":
		actionErr = existingProduct.Discontinue()
	default:
		return nil, errors.New("invalid status change action")
	}

	if actionErr != nil {
		return nil, actionErr
	}

	if err := h.repo.Update(ctx, existingProduct); err != nil {
		return nil, err
	}

	return &ChangeProductStatusResponse{
		ID:     existingProduct.ID(),
		Status: existingProduct.Status(),
	}, nil
}
