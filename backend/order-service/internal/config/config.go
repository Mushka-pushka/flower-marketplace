package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port               string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	RabbitMQURL        string
	LogLevel           string
	CatalogServiceURL  string
	PlatformCommission float64 // Процент комиссии платформы (0.0 - 1.0)
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8083"),
		DBHost:             getEnv("DB_HOST", "127.0.0.1"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "flower_user"),
		DBPassword:         getEnv("DB_PASSWORD", "flower_pass"),
		DBName:             getEnv("DB_NAME", "flower_marketplace"),
		RabbitMQURL:        getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		CatalogServiceURL:  getEnv("CATALOG_SERVICE_URL", "http://localhost:8082/api/v1/catalog"),
		PlatformCommission: getEnvAsFloat("PLATFORM_COMMISSION", 0.10),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if val, err := strconv.ParseFloat(value, 64); err == nil {
			// Ограничиваем значение от 0 до 1
			if val < 0 {
				return 0
			}
			if val > 1 {
				return 1
			}
			return val
		}
	}
	return defaultValue
}