package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/goakshit/isildur/api/handlers"
	"github.com/goakshit/isildur/platform/config"
	"github.com/goakshit/isildur/platform/constants"
	"github.com/goakshit/isildur/repositories"
	"github.com/goakshit/isildur/services"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.CFG, db *gorm.DB) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	subsRepo := repositories.NewSubscriptionsRepository(db)
	productsRepo := repositories.NewProductsRepository(db)
	subsSvc := services.NewSubscriptionService(subsRepo, productsRepo)
	productsSvc := services.NewProductsService(productsRepo)
	handler := handlers.NewHTTPHandler(subsSvc, productsSvc)

	subscriptionAPI := api.Group("/subscription")
	{
		subscriptionAPI.POST("/", handler.CreateSubscription)
		subscriptionAPI.GET(fmt.Sprintf("/:%s", constants.SubscriptionIDKey), handler.FetchSubscription)
	}
	productsAPI := api.Group("/products")
	{
		productsAPI.GET("/", handler.FetchAllProducts)
		productsAPI.GET(fmt.Sprintf("/:%s", constants.ProductIDKey), handler.Fetchproduct)
	}
	return r
}
