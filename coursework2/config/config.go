package config

import (
	"os"
)

type Config struct {
	DatabasePath string
}

func LoadConfig() *Config {
	return &Config{
		DatabasePath: getEnvOrDefault("DATABASE_PATH", "./szi_registry.db"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}