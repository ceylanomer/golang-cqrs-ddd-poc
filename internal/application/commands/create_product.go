package commands

import (
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
)

type CreateProductCommand struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type CreateProductHandler struct {
	repo product.Repository
}

func NewCreateProductHandler(repo product.Repository) *CreateProductHandler {
	return &CreateProductHandler{
		repo: repo,
	}
}

func (h *CreateProductHandler) Handle(cmd CreateProductCommand) (*product.Product, error) {
	product := &product.Product{
		ID:          uuid.New(),
		Name:        cmd.Name,
		Description: cmd.Description,
		Price:       cmd.Price,
	}

	if err := h.repo.Save(product); err != nil {
		return nil, err
	}

	return product, nil
}
