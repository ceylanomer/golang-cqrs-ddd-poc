package persistence

import (
	"fmt"
	"sync"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
)

type MemoryProductRepository struct {
	products map[uuid.UUID]*product.Product
	mutex    sync.RWMutex
}

func NewMemoryProductRepository() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[uuid.UUID]*product.Product),
	}
}

func (r *MemoryProductRepository) GetByID(id uuid.UUID) (*product.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if product, exists := r.products[id]; exists {
		return product, nil
	}
	return nil, fmt.Errorf("product not found")
}

func (r *MemoryProductRepository) Save(product *product.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.products[product.ID] = product
	return nil
}

func (r *MemoryProductRepository) Update(product *product.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return fmt.Errorf("product not found")
	}

	r.products[product.ID] = product
	return nil
}

func (r *MemoryProductRepository) Delete(id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.products[id]; !exists {
		return fmt.Errorf("product not found")
	}

	delete(r.products, id)
	return nil
}
