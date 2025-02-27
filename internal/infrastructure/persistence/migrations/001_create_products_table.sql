-- +goose Up
CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price_amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    stock_level INTEGER NOT NULL,
    stock_unit VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_price ON products(price_amount);
CREATE INDEX idx_products_name_description ON products USING gin(to_tsvector('english', name || ' ' || description));

-- +goose Down
DROP TABLE IF EXISTS products; 