package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() Config {
	return Config{
		DBHost:     env("DB_HOST", "localhost"),
		DBPort:     env("DB_PORT", "5434"),
		DBUser:     env("DB_USER", "szi_app"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     env("DB_NAME", "szi_registry"),
	}
}

func (c Config) DSN() (string, error) {
	if c.DBPassword == "" {
		return "", fmt.Errorf("переменная окружения DB_PASSWORD не задана")
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName), nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
