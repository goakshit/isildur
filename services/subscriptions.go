package services

import (
	"context"
	"time"

	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/core/ports"
	"github.com/goakshit/isildur/platform/constants"
	"github.com/google/uuid"
)

var _ ports.SubscriptionService = (*SubscriptionService)(nil)

// SubscriptionService represents required dependencies for the service.
type SubscriptionService struct {
	prodRepo ports.ProductsRepository
	subsRepo ports.SubscriptionsRepository
}

// NewSubscriptionService
func NewSubscriptionService(
	s ports.SubscriptionsRepository,
	p ports.ProductsRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subsRepo: s,
		prodRepo: p,
	}
}

func (ss SubscriptionService) CreateSubscription(
	ctx context.Context,
	pID uuid.UUID,
	durationInMonths int8,
	startDate time.Time,
) error {

	var status domain.SubscriptionStatus = domain.SubscriptionStatusInactive
	todayDate := time.Now().Truncate(24 * time.Hour)
	tomorrowDate := todayDate.Add(24 * time.Hour)
	// Check the status of subscription
	if equalDate(startDate, todayDate) {
		status = domain.SubscriptionStatusActive
	} else if startDate.After(tomorrowDate) {
		status = domain.SubscriptionStatusInactive
	} else if startDate.Before(todayDate) {
		return domain.ErrInvalidStartDate
	}

	// Fetch product details
	product, err := ss.prodRepo.GetByID(ctx, pID)
	if err != nil {
		return err
	}

	// Calculate Total cost, and tax.
	costBeforeTax := product.MonthlyPrice * float64(durationInMonths)
	taxAmount := costBeforeTax * (constants.TaxPercentApplicable / 100)
	totalCost := costBeforeTax + taxAmount

	return ss.subsRepo.Create(ctx, domain.Subscription{
		ID:               uuid.New(),
		ProductID:        product.ID,
		DurationInMonths: durationInMonths,
		Tax:              taxAmount,
		TotalCost:        totalCost,
		Status:           status,
		StartDate:        startDate,
		EndDate:          startDate.AddDate(0, int(durationInMonths), 0),
	})
}

func equalDate(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// FetchSubscription fetches subscription for a given ID.
func (ss SubscriptionService) FetchSubscription(ctx context.Context, id uuid.UUID) (domain.Subscription, error) {
	if id == uuid.Nil {
		return domain.Subscription{}, domain.ErrSubscriptionIDIsInvalid
	}
	return ss.subsRepo.GetByID(ctx, id)
}

func (ss SubscriptionService) UpdateSubscriptionStatus(ctx context.Context, id uuid.UUID, status domain.SubscriptionStatus) error {
	if id == uuid.Nil {
		return domain.ErrSubscriptionIDIsInvalid
	}
	return ss.subsRepo.Patch(ctx, id, map[string]interface{}{
		"status": status,
	})
}
