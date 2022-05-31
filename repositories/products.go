package repositories

import (
	"context"
	"errors"

	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ ports.ProductsRepository = (*ProductsRepository)(nil)

// ProductsRepository represents list of dependencies for repository.
type ProductsRepository struct {
	db *gorm.DB
}

// NewProductsRepository creates and returns new ProductsRepository.
func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

// GetByID returns product by id from db.
func (cr ProductsRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	var product domain.Product
	result := cr.db.WithContext(ctx).Where(domain.Product{
		ID: id,
	}).First(&product)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return product, domain.ErrProductNotfound
	}
	return product, result.Error
}

// GetAll fetches all the products in the database.
func (cr ProductsRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	result := cr.db.WithContext(ctx).Find(&products)
	return products, result.Error
}
