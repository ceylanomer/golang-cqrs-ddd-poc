package product

import (
	"context"

	"github.com/google/uuid"
)

// Repository represents the main write repository for products
type Repository interface {
	Save(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
}

// ReadOnlyRepository represents a read-only repository for product queries
type ReadOnlyRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*ProductReadModel, error)
	FindAll(ctx context.Context, filter ProductFilter) ([]ProductReadModel, error)
	FindByStatus(ctx context.Context, status ProductStatus) ([]ProductReadModel, error)
}

// ProductReadModel represents a denormalized view of the Product aggregate
type ProductReadModel struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	PriceAmount float64       `json:"price_amount"`
	Currency    string        `json:"currency"`
	StockLevel  int           `json:"stock_level"`
	StockUnit   string        `json:"stock_unit"`
	Status      ProductStatus `json:"status"`
	Version     int           `json:"version"`
}

// ProductFilter represents query filters for products
type ProductFilter struct {
	MinPrice   *float64
	MaxPrice   *float64
	Status     *ProductStatus
	StockLevel *int
	SearchTerm string
	PageSize   int
	PageNumber int
}


