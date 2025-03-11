package queries

import (
	"context"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/gofiber/fiber/v2"
)

type ListProductsQuery struct {
	MinPrice   *float64               `query:"min_price"`
	MaxPrice   *float64               `query:"max_price"`
	Status     *product.ProductStatus `query:"status"`
	StockLevel *int                   `query:"stock_level"`
	SearchTerm string                 `query:"search"`
	PageSize   int                    `query:"page_size"`
	PageNumber int                    `query:"page"`
}

type ListProductsResponse struct {
	Products []product.ProductReadModel `json:"products"`
	Total    int                        `json:"total"`
}

type ListProductsHandler struct {
	repo product.ReadOnlyRepository
}

func NewListProductsHandler(repo product.ReadOnlyRepository) *ListProductsHandler {
	return &ListProductsHandler{repo: repo}
}

func (h *ListProductsHandler) Handle(ctx context.Context, query *ListProductsQuery) (*ListProductsResponse, error) {
	filter := product.ProductFilter{
		MinPrice:   query.MinPrice,
		MaxPrice:   query.MaxPrice,
		Status:     query.Status,
		StockLevel: query.StockLevel,
		SearchTerm: query.SearchTerm,
		PageSize:   query.PageSize,
		PageNumber: query.PageNumber,
	}

	products, err := h.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &ListProductsResponse{
		Products: products,
		Total:    len(products),
	}, nil
}
