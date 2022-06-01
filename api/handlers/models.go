package handlers

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type CreateSubscriptionRequest struct {
	ProductID        string `json:"product_id" valid:"required,uuidv4"`
	StartDate        string `json:"start_date" valid:"required"`
	DurationInMonths int8   `json:"duration_in_months" valid:"required,numeric"`
}
