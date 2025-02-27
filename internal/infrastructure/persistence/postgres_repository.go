package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Write Repository Implementation
func (r *PostgresRepository) Save(ctx context.Context, product *product.Product) error {
	query := `
		INSERT INTO products (
			id, name, description, price_amount, currency,
			stock_level, stock_unit, status, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		product.ID(),
		product.Name(),
		product.Description(),
		product.Price(),
		product.Price(),
		product.Stock(),
		product.Stock(),
		product.Status(),
		product.Version(),
	)

	return err
}

func (r *PostgresRepository) Update(ctx context.Context, product *product.Product) error {
	query := `
		UPDATE products SET
			name = $2,
			description = $3,
			price_amount = $4,
			currency = $5,
			stock_level = $6,
			stock_unit = $7,
			status = $8,
			version = $9
		WHERE id = $1 AND version = $10
	`

	result, err := r.db.ExecContext(ctx, query,
		product.ID(),
		product.Name(),
		product.Description(),
		product.Price(),
		product.Price(),
		product.Stock(),
		product.Stock(),
		product.Status(),
		product.Version()+1, // Increment version
		product.Version(),   // Current version for optimistic locking
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product version mismatch or not found")
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	query := `
		SELECT id, name, description, price_amount, currency,
			   stock_level, stock_unit, status, version
		FROM products
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var (
		productID   uuid.UUID
		name        string
		description string
		priceAmount float64
		currency    string
		stockLevel  int
		stockUnit   string
		status      product.ProductStatus
		version     int
	)

	if err := row.Scan(
		&productID, &name, &description, &priceAmount, &currency,
		&stockLevel, &stockUnit, &status, &version,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	price, err := product.NewPrice(priceAmount, currency)
	if err != nil {
		return nil, err
	}

	stock, err := product.NewStock(stockLevel, stockUnit)
	if err != nil {
		return nil, err
	}

	// Note: This is a simplified version. In a real implementation,
	// you might need to reconstruct the product with all its internal state
	return product.NewProduct(
		name, description, price, stock,
	)
}

// Read Repository Implementation
func (r *PostgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*product.ProductReadModel, error) {
	query := `
		SELECT id, name, description, price_amount, currency,
			   stock_level, stock_unit, status, version
		FROM products
		WHERE id = $1
	`

	var readModel product.ProductReadModel
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&readModel.ID, &readModel.Name, &readModel.Description,
		&readModel.PriceAmount, &readModel.Currency,
		&readModel.StockLevel, &readModel.StockUnit,
		&readModel.Status, &readModel.Version,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}

	return &readModel, err
}

func (r *PostgresRepository) FindAll(ctx context.Context, filter product.ProductFilter) ([]product.ProductReadModel, error) {
	query := `
		SELECT id, name, description, price_amount, currency,
			   stock_level, stock_unit, status, version
		FROM products
		WHERE ($1::float8 IS NULL OR price_amount >= $1)
		AND ($2::float8 IS NULL OR price_amount <= $2)
		AND ($3::text IS NULL OR status = $3)
		AND ($4::int IS NULL OR stock_level >= $4)
		AND ($5::text IS NULL OR 
			 name ILIKE '%' || $5 || '%' OR 
			 description ILIKE '%' || $5 || '%')
		LIMIT $6 OFFSET $7
	`

	rows, err := r.db.QueryContext(ctx, query,
		filter.MinPrice,
		filter.MaxPrice,
		filter.Status,
		filter.StockLevel,
		filter.SearchTerm,
		filter.PageSize,
		filter.PageSize*filter.PageNumber,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.ProductReadModel
	for rows.Next() {
		var p product.ProductReadModel
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description,
			&p.PriceAmount, &p.Currency,
			&p.StockLevel, &p.StockUnit,
			&p.Status, &p.Version,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, rows.Err()
}

func (r *PostgresRepository) FindByStatus(ctx context.Context, status product.ProductStatus) ([]product.ProductReadModel, error) {
	query := `
		SELECT id, name, description, price_amount, currency,
			   stock_level, stock_unit, status, version
		FROM products
		WHERE status = $1
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.ProductReadModel
	for rows.Next() {
		var p product.ProductReadModel
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description,
			&p.PriceAmount, &p.Currency,
			&p.StockLevel, &p.StockUnit,
			&p.Status, &p.Version,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, rows.Err()
}
