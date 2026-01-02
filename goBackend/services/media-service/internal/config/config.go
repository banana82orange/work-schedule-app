package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	GRPCPort    int
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	StoragePath string
	StorageURL  string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		GRPCPort:    getEnvInt("GRPC_PORT", 50055),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnvInt("DB_PORT", 5432),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "portfolio"),
		DBSSLMode:   getEnv("DB_SSL_MODE", "disable"),
		StoragePath: getEnv("STORAGE_PATH", "./uploads"),
		StorageURL:  getEnv("STORAGE_URL", "http://localhost:50055/files"),
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
