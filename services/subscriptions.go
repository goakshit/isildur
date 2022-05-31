package services

import (
	"github.com/goakshit/isildur/core/ports"
)

var _ ports.SubscriptionService = (*SubscriptionService)(nil)

// SubscriptionService represents required dependencies for the service.
type SubscriptionService struct {
	subsRepo ports.SubscriptionsRepository
}

// NewSubscriptionService
func NewSubscriptionService(
	s ports.SubscriptionsRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subsRepo: s,
	}
}
