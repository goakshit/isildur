package config

import (
	"os"
)

// CFG represents root structure of env configuration of the service.
type CFG struct {
	ServiceName  string
	ServicePort  string
	ServiceLevel string
	DB           DBConfig
}

// DBConfig represents configuration used to connect with the db.
type DBConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

// LoadFromEnv will load the env vars from the OS.
func LoadFromEnv() *CFG {
	return &CFG{
		ServiceName:  getEnv("SERVICE_NAME", "subscription-service"),
		ServicePort:  getEnv("SERVICE_PORT", "8080"),
		ServiceLevel: getEnv("SERVICE_LEVEL", "debug"),
		DB: DBConfig{
			User: getEnv("DB_USER", "gymondo_user"),
			Pass: getEnv("DB_PASS", "gymondo_pass"),
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "5432"),
			Name: getEnv("DB_NAME", "gymondo"),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
