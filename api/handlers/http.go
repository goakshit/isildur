package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/core/ports"
	"github.com/goakshit/isildur/platform/constants"
	"github.com/google/uuid"
)

// HTTPHandler holds dependencies used inside the http handlers.
type HTTPHandler struct {
	Subs     ports.SubscriptionService
	Products ports.ProductsService
}

// NewHTTPHandler returns a new HTTPHandler.
func NewHTTPHandler(
	subs ports.SubscriptionService,
	products ports.ProductsService,
) HTTPHandler {
	return HTTPHandler{
		Subs:     subs,
		Products: products,
	}
}

// CreateSubscription creates a subscription.
func (h *HTTPHandler) CreateSubscription(ctx *gin.Context) {
	r := CreateSubscriptionRequest{}
	if err := ctx.BindJSON(&r); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	pid, err := uuid.Parse(r.ProductID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	// Parsing the start date.
	date, err := time.Parse(constants.DateFormat, r.StartDate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      domain.ErrInvalidStartDate.Error(),
		})
		return
	}

	if err = h.Subs.CreateSubscription(ctx, pid, r.DurationInMonths, date); err != nil {
		errResp := mapErrorResponseFromError(err)
		ctx.AbortWithStatusJSON(errResp.StatusCode, errResp)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"message":     "Successfully create subscription",
	})
}

// FetchAllProducts fetches all the available products.
func (h *HTTPHandler) FetchAllProducts(ctx *gin.Context) {
	products, err := h.Products.FetchAllProducts(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

// FetchProduct fetches all the available products.
func (h *HTTPHandler) FetchProduct(ctx *gin.Context) {
	pID, err := uuid.Parse(ctx.Param(constants.ProductIDKey))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}
	product, err := h.Products.FetchProduct(ctx, pID)
	if err != nil {
		errResp := mapErrorResponseFromError(err)
		ctx.AbortWithStatusJSON(errResp.StatusCode, errResp)
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// FetchSubscription fetches subscription details for given id.
func (h *HTTPHandler) FetchSubscription(ctx *gin.Context) {
	sID, err := uuid.Parse(ctx.Param(constants.SubscriptionIDKey))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}
	subscription, err := h.Subs.FetchSubscription(ctx, sID)
	if err != nil {
		errResp := mapErrorResponseFromError(err)
		ctx.AbortWithStatusJSON(errResp.StatusCode, errResp)
		return
	}
	ctx.JSON(http.StatusOK, subscription)
}

// UpdateSubscriptionStatus updates the state of subscription for a given id.
func (h *HTTPHandler) UpdateSubscriptionStatus(ctx *gin.Context) {
	sID, err := uuid.Parse(ctx.Param(constants.SubscriptionIDKey))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	updateStatus := ctx.Query("status")
	if updateStatus == "" || domain.MapStringToSubscriptionStatus(updateStatus) == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      domain.ErrInvalidSubscriptionStatusPassed.Error(),
		})
		return
	}

	subscription, err := h.Subs.FetchSubscription(ctx, sID)
	if err != nil {
		errResp := mapErrorResponseFromError(err)
		ctx.AbortWithStatusJSON(errResp.StatusCode, errResp)
		return
	}

	// Error handling: Don't allow subscription to be active if it is cancelled
	if subscription.Status == domain.SubscriptionStatusCancel {
		errResp := mapErrorResponseFromError(domain.ErrCannotUpdateCancelledSubscription)
		ctx.AbortWithStatusJSON(errResp.StatusCode, errResp)
		return
	}

	if err = h.Subs.UpdateSubscriptionStatus(ctx, sID, domain.MapStringToSubscriptionStatus(updateStatus)); err != nil {
		errResp := mapErrorResponseFromError(err)
		ctx.AbortWithStatusJSON(errResp.StatusCode, errResp)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"message":     "Successfully updated the subscription status.",
	})
}

func mapErrorResponseFromError(err error) ErrorResponse {
	resp := ErrorResponse{
		Error:      err.Error(),
		StatusCode: http.StatusInternalServerError,
	}
	if errors.Is(err, domain.ErrProductNotfound) ||
		errors.Is(err, domain.ErrSubscriptionNotfound) {

		resp.StatusCode = http.StatusNotFound

	} else if errors.Is(err, domain.ErrProductIDIsInvalid) ||
		errors.Is(err, domain.ErrInvalidStartDate) ||
		errors.Is(err, domain.ErrSubscriptionIDIsInvalid) ||
		errors.Is(err, domain.ErrCannotUpdateCancelledSubscription) {

		resp.StatusCode = http.StatusBadRequest

	}
	return resp
}
