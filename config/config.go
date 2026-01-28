package config

import (
	"fmt"
	"os"
)

// содержит настройки приложения
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "5433"),
		DBUser:     getEnvOrDefault("DB_USER", "szi_user"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "MyV3ryS3cur3P@ss2026!"),
		DBName:     getEnvOrDefault("DB_NAME", "szi_registry"),
	}
}

// возвращает значение переменной окружения или значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// возвращает строку подключения к базе данных
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
