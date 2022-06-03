package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/goakshit/isildur/platform/config"
	"github.com/goakshit/isildur/platform/constants"
	"github.com/goakshit/isildur/repositories"
	"github.com/goakshit/isildur/services"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, cfg *config.CFG, db *gorm.DB) {
	api := r.Group("/api")

	subsRepo := repositories.NewSubscriptionsRepository(db)
	productsRepo := repositories.NewProductsRepository(db)
	subsSvc := services.NewSubscriptionService(subsRepo, productsRepo)
	productsSvc := services.NewProductsService(productsRepo)
	handler := NewHTTPHandler(subsSvc, productsSvc)

	subscriptionAPI := api.Group("/subscription")
	{
		subscriptionAPI.POST("/", handler.CreateSubscription)
		subscriptionAPI.GET(fmt.Sprintf("/:%s", constants.SubscriptionIDKey), handler.FetchSubscription)
		subscriptionAPI.PATCH(fmt.Sprintf("/:%s", constants.SubscriptionIDKey), handler.UpdateSubscriptionStatus)
	}
	productsAPI := api.Group("/products")
	{
		productsAPI.GET("/", handler.FetchAllProducts)
		productsAPI.GET(fmt.Sprintf("/:%s", constants.ProductIDKey), handler.FetchProduct)
	}
}
