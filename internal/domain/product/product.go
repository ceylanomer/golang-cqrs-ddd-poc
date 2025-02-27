package product

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Value Objects
type Price struct {
	amount   float64
	currency string
}

func NewPrice(amount float64, currency string) (Price, error) {
	if amount < 0 {
		return Price{}, errors.New("price cannot be negative")
	}
	if currency == "" {
		return Price{}, errors.New("currency is required")
	}
	return Price{amount: amount, currency: currency}, nil
}

type Stock struct {
	quantity int
	unit     string
}

func NewStock(quantity int, unit string) (Stock, error) {
	if quantity < 0 {
		return Stock{}, errors.New("stock quantity cannot be negative")
	}
	if unit == "" {
		return Stock{}, errors.New("unit is required")
	}
	return Stock{quantity: quantity, unit: unit}, nil
}

// Product Entity (Aggregate Root)
type Product struct {
	id          uuid.UUID
	name        string
	description string
	price       Price
	stock       Stock
	status      ProductStatus
	version     int
	events      []ProductEvent
}

type ProductStatus string

const (
	StatusDraft        ProductStatus = "DRAFT"
	StatusActive       ProductStatus = "ACTIVE"
	StatusInactive     ProductStatus = "INACTIVE"
	StatusDiscontinued ProductStatus = "DISCONTINUED"
)

// Factory method
func NewProduct(name, description string, price Price, stock Stock) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name is required")
	}

	return &Product{
		id:          uuid.New(),
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		status:      StatusDraft,
		version:     1,
	}, nil
}

// Business methods
func (p *Product) Activate() error {
	if p.status == StatusDiscontinued {
		return errors.New("cannot activate discontinued product")
	}
	if p.stock.quantity == 0 {
		return errors.New("cannot activate product with zero stock")
	}
	p.status = StatusActive
	p.version++
	return nil
}

func (p *Product) Deactivate() error {
	if p.status == StatusDiscontinued {
		return errors.New("cannot deactivate discontinued product")
	}
	p.status = StatusInactive
	p.version++
	return nil
}

func (p *Product) UpdatePrice(newPrice Price) error {
	if p.status == StatusDiscontinued {
		return errors.New("cannot update price of discontinued product")
	}
	p.price = newPrice
	p.version++
	return nil
}

func (p *Product) UpdateStock(quantity int) error {
	if quantity < 0 {
		return errors.New("stock quantity cannot be negative")
	}
	if p.status == StatusDiscontinued {
		return errors.New("cannot update stock of discontinued product")
	}

	newStock, err := NewStock(quantity, p.stock.unit)
	if err != nil {
		return err
	}

	p.stock = newStock
	p.version++
	return nil
}

func (p *Product) Discontinue() error {
	if p.status == StatusDiscontinued {
		return errors.New("product is already discontinued")
	}
	p.status = StatusDiscontinued
	p.version++
	return nil
}

// Getters (since fields are private)
func (p *Product) ID() uuid.UUID         { return p.id }
func (p *Product) Name() string          { return p.name }
func (p *Product) Description() string   { return p.description }
func (p *Product) Price() float64        { return p.price.amount }
func (p *Product) Currency() string      { return p.price.currency }
func (p *Product) Stock() int            { return p.stock.quantity }
func (p *Product) StockUnit() string     { return p.stock.unit }
func (p *Product) Status() ProductStatus { return p.status }
func (p *Product) Version() int          { return p.version }

// Price methods
func (p Price) Amount() float64 {
	return p.amount
}

func (p Price) Currency() string {
	return p.currency
}

// Stock methods
func (s Stock) Quantity() int {
	return s.quantity
}

func (s Stock) Unit() string {
	return s.unit
}

type ProductEvent interface {
	EventName() string
}

type ProductPriceChanged struct {
	ProductID uuid.UUID
	OldPrice  Price
	NewPrice  Price
	ChangedAt time.Time
}
