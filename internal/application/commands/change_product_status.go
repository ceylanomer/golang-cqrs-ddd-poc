package commands

import (
	"context"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
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
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	
	if existingProduct.Version() != cmd.Version {
		return nil, fiber.NewError(fiber.StatusConflict, "product has been modified by another process")
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
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid status change action")
	}

	if actionErr != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, actionErr.Error())
	}

	if err := h.repo.Update(ctx, existingProduct); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &ChangeProductStatusResponse{
		ID:     existingProduct.ID(),
		Status: existingProduct.Status(),
	}, nil
}
