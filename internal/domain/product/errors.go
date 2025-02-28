package product

import "errors"

var (
	ErrNotFound        = errors.New("product not found")
	ErrVersionMismatch = errors.New("product version mismatch")
)
