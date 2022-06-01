package repositories

import (
	"context"
	"errors"

	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ ports.SubscriptionsRepository = (*SubscriptionsRepository)(nil)

// SubscriptionsRepository represents list of dependencies for repository.
type SubscriptionsRepository struct {
	db *gorm.DB
}

// NewSubscriptionsRepository creates and returns new SubscriptionsRepository.
func NewSubscriptionsRepository(db *gorm.DB) *SubscriptionsRepository {
	return &SubscriptionsRepository{
		db: db,
	}
}

// GetByID fetches subscription for a given id.
func (sr SubscriptionsRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	subscription := domain.Subscription{}
	result := sr.db.WithContext(ctx).Where(domain.Subscription{
		ID: id,
	}).First(&subscription)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return subscription, domain.ErrSubscriptionNotfound
	}
	return subscription, result.Error
}

// Create is used to create a subscription in the db.
func (sr SubscriptionsRepository) Create(ctx context.Context, sub domain.Subscription) error {
	return sr.db.WithContext(ctx).Create(&sub).Error
}

// Patch updates the data in subscription for a given id.
func (sr SubscriptionsRepository) Patch(ctx context.Context, id uuid.UUID, update map[string]interface{}) error {
	updateOP := sr.db.WithContext(ctx).Model(&domain.Subscription{}).
		Where(&domain.Subscription{
			ID: id,
		}).
		Updates(update)
	if updateOP.Error != nil {
		return updateOP.Error
	}

	// Incorrect ID
	if updateOP.RowsAffected == 0 {
		return domain.ErrSubscriptionNotfound
	}
	return nil
}
