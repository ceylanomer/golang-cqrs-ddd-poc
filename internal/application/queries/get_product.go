package queries

import (
	"context"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
)

type GetProductQuery struct {
	ID uuid.UUID `params:"id"`
}

type GetProductHandler struct {
	repo product.ReadOnlyRepository
}

func NewGetProductHandler(repo product.ReadOnlyRepository) *GetProductHandler {
	return &GetProductHandler{repo: repo}
}

func (h *GetProductHandler) Handle(ctx context.Context, query *GetProductQuery) (*product.ProductReadModel, error) {
	return h.repo.FindByID(ctx, query.ID)
}
