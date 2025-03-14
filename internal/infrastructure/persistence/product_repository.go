package persistence

import (
	"context"

	"github.com/ceylanomer/golang-cqrs-ddd-poc/internal/domain/product"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Write Repository Implementation
func (r *ProductRepository) Save(ctx context.Context, product *product.Product) error {
	model := FromDomain(product)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ProductRepository) Update(ctx context.Context, product *product.Product) error {
	model := FromDomain(product)
	result := r.db.WithContext(ctx).Model(&ProductModel{}).
		Where("id = ? AND version = ?", model.ID, model.Version-1).
		Updates(model)

	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusConflict, "product has been modified by another process")
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&ProductModel{}, id)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "product not found")
	}

	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	var model ProductModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "product not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return model.ToDomain()
}

// Read Repository Implementation
func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*product.ProductReadModel, error) {
	var model ProductModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "product not found")
		}
		return nil, err
	}

	return &product.ProductReadModel{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		PriceAmount: model.PriceAmount,
		Currency:    model.Currency,
		StockLevel:  model.StockLevel,
		StockUnit:   model.StockUnit,
		Status:      model.Status,
		Version:     model.Version,
	}, nil
}

func (r *ProductRepository) FindAll(ctx context.Context, filter product.ProductFilter) ([]product.ProductReadModel, error) {
	var models []ProductModel
	query := r.db.WithContext(ctx)

	if filter.MinPrice != nil {
		query = query.Where("price_amount >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price_amount <= ?", *filter.MaxPrice)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.StockLevel != nil {
		query = query.Where("stock_level >= ?", *filter.StockLevel)
	}
	if filter.SearchTerm != "" {
		query = query.Where(
			"name ILIKE ? OR description ILIKE ?",
			"%"+filter.SearchTerm+"%",
			"%"+filter.SearchTerm+"%",
		)
	}

	// Apply pagination
	if filter.PageSize > 0 {
		offset := filter.PageSize * filter.PageNumber
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Convert to read models
	readModels := make([]product.ProductReadModel, len(models))
	for i, model := range models {
		readModels[i] = product.ProductReadModel{
			ID:          model.ID,
			Name:        model.Name,
			Description: model.Description,
			PriceAmount: model.PriceAmount,
			Currency:    model.Currency,
			StockLevel:  model.StockLevel,
			StockUnit:   model.StockUnit,
			Status:      model.Status,
			Version:     model.Version,
		}
	}

	return readModels, nil
}

func (r *ProductRepository) FindByStatus(ctx context.Context, status product.ProductStatus) ([]product.ProductReadModel, error) {
	var models []ProductModel
	if err := r.db.WithContext(ctx).Where("status = ?", status).Find(&models).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Convert to read models
	readModels := make([]product.ProductReadModel, len(models))
	for i, model := range models {
		readModels[i] = product.ProductReadModel{
			ID:          model.ID,
			Name:        model.Name,
			Description: model.Description,
			PriceAmount: model.PriceAmount,
			Currency:    model.Currency,
			StockLevel:  model.StockLevel,
			StockUnit:   model.StockUnit,
			Status:      model.Status,
			Version:     model.Version,
		}
	}

	return readModels, nil
}
