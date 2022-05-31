package ports

//go:generate mockgen -package=ports -source=ports.go -destination=mocks.go

import (
	"context"

	"github.com/goakshit/isildur/core/domain"
	"github.com/google/uuid"
)

// ProductsRepository describers database operations on products entity.
type ProductsRepository interface {
	// GetAll fetches all the products in the database.
	GetAll(ctx context.Context) ([]domain.Product, error)
	// GetByID fetches product for a given id.
	GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
}

// SubscriptionsRepository describers database operations on subscriptions entity.
type SubscriptionsRepository interface {
	// Create is used to create a subscription in the db.
	Create(ctx context.Context, sub domain.Subscription) error
	// GetByID fetches subscription for a given id.
	GetByID(ctx context.Context, id uuid.UUID) (domain.Subscription, error)
	// Patch updates the data in subscription for a given id.
	Patch(ctx context.Context, id uuid.UUID, update map[string]interface{}) error
}

// SubscriptionService describes main business functionality of subscription service.
type SubscriptionService interface {
	//
	// CreateSubscription(ctx context.Context, )
}

// ProductsService describes main business functionality of products.
type ProductsService interface {
	// FetchAllProduct fetches all the products in the database.
	FetchAllProducts(ctx context.Context) ([]domain.Product, error)
	// FetchProduct fetches product for a given ID.
	FetchProduct(ctx context.Context, id uuid.UUID) (domain.Product, error)
}
