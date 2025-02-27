package product

import (
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
}

type Repository interface {
	GetByID(id uuid.UUID) (*Product, error)
	Save(product *Product) error
	Update(product *Product) error
	Delete(id uuid.UUID) error
}
