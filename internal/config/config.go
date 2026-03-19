package config

import (
	"os"
)

// Config represents the application configuration.
type Config struct {
	JWTSecret string
	JWTIssuer string
}

// LoadConfig loads the application configuration from environment variables.
func LoadConfig() *Config {
	return &Config{
		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTIssuer: getEnv("JWT_ISSUER", ""),
	}
}

// getEnv retrieves the value of an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
