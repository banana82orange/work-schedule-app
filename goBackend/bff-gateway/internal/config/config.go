package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the BFF Gateway configuration
type Config struct {
	// HTTP Server
	HTTPPort int

	// Service URLs
	AuthServiceURL      string
	ProjectServiceURL   string
	TaskServiceURL      string
	AnalyticsServiceURL string
	MediaServiceURL     string

	// JWT
	JWTSecret string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load environment variables")
	}
	return &Config{
		HTTPPort:            getEnvInt("HTTP_PORT", 8080),
		AuthServiceURL:      getEnv("AUTH_SERVICE_URL", "localhost:50051"),
		ProjectServiceURL:   getEnv("PROJECT_SERVICE_URL", "localhost:50052"),
		TaskServiceURL:      getEnv("TASK_SERVICE_URL", "localhost:50053"),
		AnalyticsServiceURL: getEnv("ANALYTICS_SERVICE_URL", "localhost:50054"),
		MediaServiceURL:     getEnv("MEDIA_SERVICE_URL", "localhost:50055"),
		JWTSecret:           getEnv("JWT_SECRET", "development-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
