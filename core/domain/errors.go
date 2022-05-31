package domain

import "errors"

var (
	// ErrProductNotfound is the error used when a product doesn't exist for a given id.
	ErrProductNotfound = errors.New("product not found")

	// ErrSubscriptionNotfound is the error used when a subscriptuon doesn't exist for a given id.
	ErrSubscriptionNotfound = errors.New("subscription not found")

	// ErrProductIDIsInvalid is the error used when a given product id is an invalid uuid.
	ErrProductIDIsInvalid = errors.New("invalid product id")
)
