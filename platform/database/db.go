package database

import (
	"fmt"

	"github.com/goakshit/isildur/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Singleton
var db *gorm.DB

func getGORMConfig() *gorm.Config {
	return &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// Doesn't pluralize the table names
			// Eg: 'user' table won't be pluralized to 'users' table
			SingularTable: true,
		},
	}
}

// Return postgres connection string
func getPostgresConnString(cfg *config.CFG) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
}

// GetGormClient - Returns db client
func GetGormClient(cfg *config.CFG) *gorm.DB {
	var err error
	db, err = gorm.Open(postgres.Open(getPostgresConnString(cfg)), getGORMConfig())
	if err != nil {
		panic("Failed to open postgres connection\n" + err.Error())
	}
	return db
}
