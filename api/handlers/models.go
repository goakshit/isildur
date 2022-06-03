package handlers

// ErrorResponse represents the standard error response that gets sent
// for every non 200 request.
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// CreateSubscriptionRequest represents the request structure for create
// subscription endpoint.
type CreateSubscriptionRequest struct {
	ProductID        string `json:"product_id" valid:"required,uuidv4"`
	StartDate        string `json:"start_date" valid:"required"`
	DurationInMonths int8   `json:"duration_in_months" valid:"required,numeric"`
}
