package domain

import (
	"time"

	"github.com/google/uuid"
)

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
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	ProductID        uuid.UUID `json:"-"`
	DurationInMonths int8      `json:"duration_in_months"`
	Tax              float64   `json:"tax"`
	TotalCost        float64   `json:"total_cost"`
	Status           string    `json:"status"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
}
