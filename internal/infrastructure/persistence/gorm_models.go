package persistence

import (
	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductModel is the GORM model for Product aggregate
type ProductModel struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"not null"`
	Description string
	PriceAmount float64               `gorm:"not null"`
	Currency    string                `gorm:"not null;size:3"`
	StockLevel  int                   `gorm:"not null"`
	StockUnit   string                `gorm:"not null"`
	Status      product.ProductStatus `gorm:"not null"`
	Version     int                   `gorm:"not null"`
}

// TableName overrides the table name
func (ProductModel) TableName() string {
	return "products"
}

// ToDomain converts GORM model to domain model
func (p *ProductModel) ToDomain() (*product.Product, error) {
	price, err := product.NewPrice(p.PriceAmount, p.Currency)
	if err != nil {
		return nil, err
	}

	stock, err := product.NewStock(p.StockLevel, p.StockUnit)
	if err != nil {
		return nil, err
	}

	// Since we can't directly create a product with an existing ID,
	// we'll need to create it first and then set the fields
	prod, err := product.NewProduct(p.Name, p.Description, price, stock)
	if err != nil {
		return nil, err
	}

	// Set the remaining fields that were loaded from the database
	// Note: This assumes we have access to these setters in the domain model
	// You might need to add these methods to your domain model
	prod.SetID(p.ID)
	prod.SetStatus(p.Status)
	prod.SetVersion(p.Version)

	return prod, nil
}

// FromDomain creates a GORM model from domain model
func FromDomain(p *product.Product) *ProductModel {
	return &ProductModel{
		ID:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		PriceAmount: p.Price(),
		Currency:    p.Currency(),
		StockLevel:  p.Stock(),
		StockUnit:   p.StockUnit(),
		Status:      p.Status(),
		Version:     p.Version(),
	}
}
