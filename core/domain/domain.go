// Domain package is responsible for all domain related objects,
// like errors, subscription statuses and so on.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionStatus reprsents the applicable status for the a subscription.
type SubscriptionStatus string

func (s SubscriptionStatus) String() string {
	return string(s)
}

var (

	// SubscriptionStatusActive represents activate subscription status
	SubscriptionStatusActive SubscriptionStatus = "active"

	// SubscriptionStatusPaused represents paused subscription status
	SubscriptionStatusPaused SubscriptionStatus = "paused"

	// SubscriptionStatusCancel represents cancel subscription status
	SubscriptionStatusCancel SubscriptionStatus = "cancelled"

	// SubscriptionStatusInactive represents inactive subscription status
	SubscriptionStatusInactive SubscriptionStatus = "inactive"
)

// MapStringToSubscriptionStatus maps string literal to SubscriptionStatus type.
func MapStringToSubscriptionStatus(status string) SubscriptionStatus {
	switch status {
	case "active":
		return SubscriptionStatusActive
	case "paused":
		return SubscriptionStatusPaused
	case "cancelled":
		return SubscriptionStatusCancel
	case "inactive":
		return SubscriptionStatusInactive
	default:
		return ""
	}
}

// Product represents structure for product entity in db.
type Product struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	MonthlyPrice   float64   `json:"monthly_price"`
	InstructorName string    `json:"instructor_name"` // This should ideally be fk to instructor or users table
}

// Subscription represents structure for subscription entity in db.
type Subscription struct {
	ID               uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;"`
	ProductID        uuid.UUID          `json:"-"`
	DurationInMonths int8               `json:"duration_in_months"`
	Tax              float64            `json:"tax"`
	TotalCost        float64            `json:"total_cost"`
	Status           SubscriptionStatus `json:"status"`
	StartDate        time.Time          `json:"start_date"`
	EndDate          time.Time          `json:"end_date"`
}
