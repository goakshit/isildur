// Package
package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/goakshit/isildur/api/router"
	"github.com/goakshit/isildur/platform/config"
	"github.com/goakshit/isildur/platform/database"
)

func main() {

	cfg := config.LoadFromEnv()
	db := database.GetGormClient(cfg)

	// Set gin mode in different environment
	gin.SetMode(cfg.ServiceLevel)
	r := router.SetupRouter(cfg, db)
	if err := r.Run(fmt.Sprintf(":%s", cfg.ServicePort)); err != nil {
		log.Fatalln("failed to setup router")
	}
}
