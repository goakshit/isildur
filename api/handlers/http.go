package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/core/ports"
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
	// var request marqetagen.CardProductRequest
	// if err := ctx.BindJSON(&request); err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
	// 	return
	// }
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
func (h *HTTPHandler) Fetchproduct(ctx *gin.Context) {
	pID, err := uuid.Parse(ctx.Param("product-id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}
	product, err := h.Products.FetchProduct(ctx, pID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotfound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
				StatusCode: http.StatusNotFound,
				Error:      err.Error(),
			})
		} else if errors.Is(err, domain.ErrProductIDIsInvalid) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err.Error(),
			})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, product)
}
