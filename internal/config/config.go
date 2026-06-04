package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	ServerPort string
	GinMode    string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "34501"),
		GinMode:    getEnv("GIN_MODE", "debug"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
