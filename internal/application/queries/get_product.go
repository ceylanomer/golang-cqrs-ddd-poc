package queries

import (
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
)

type GetProductQuery struct {
	ID uuid.UUID
}

type GetProductHandler struct {
	repo product.Repository
}

func NewGetProductHandler(repo product.Repository) *GetProductHandler {
	return &GetProductHandler{
		repo: repo,
	}
}

func (h *GetProductHandler) Handle(query GetProductQuery) (*product.Product, error) {
	return h.repo.GetByID(query.ID)
}
