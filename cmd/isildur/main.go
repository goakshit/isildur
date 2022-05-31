// Package
package main

import (
	"log"

	"github.com/goakshit/isildur/api/router"
	"github.com/goakshit/isildur/core/domain"
	"github.com/goakshit/isildur/platform/config"
	"github.com/goakshit/isildur/platform/database"
)

func main() {

	cfg := config.LoadFromEnv()
	db := database.GetGormClient(cfg)

	db.AutoMigrate(&domain.Product{}, &domain.Subscription{})

	r := router.SetupRouter(cfg, db)
	if err := r.Run(); err != nil {
		log.Fatalln("failed to setup router")
	}
}
