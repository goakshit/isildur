package services

import (
	"context"

	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/core/ports"
	"github.com/google/uuid"
)

var _ ports.ProductsService = (*ProductsService)(nil)

// ProductsService represents required dependencies for the service.
type ProductsService struct {
	ProductsRepo ports.ProductsRepository
}

// NewSubscriptionService
func NewProductsService(
	p ports.ProductsRepository,
) *ProductsService {
	return &ProductsService{
		ProductsRepo: p,
	}
}

// FetchAllProduct fetches all the products in the database.
func (p ProductsService) FetchAllProducts(ctx context.Context) ([]domain.Product, error) {
	return p.ProductsRepo.GetAll(ctx)
}

// FetchProduct fetches product for a given id.
func (p ProductsService) FetchProduct(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	if id == uuid.Nil {
		return domain.Product{}, domain.ErrProductIDIsInvalid
	}
	return p.ProductsRepo.GetByID(ctx, id)
}
