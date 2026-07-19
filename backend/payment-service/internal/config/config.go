package config

import (
	"os"
)

type Config struct {
	Port            string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	RabbitMQURL     string
	OrderServiceURL string 
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8084"),
		DBHost:          getEnv("DB_HOST", "127.0.0.1"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "flower_user"),
		DBPassword:      getEnv("DB_PASSWORD", "flower_pass"),
		DBName:          getEnv("DB_NAME", "flower_marketplace"),
		RabbitMQURL:     getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		OrderServiceURL: getEnv("ORDER_SERVICE_URL", "http://localhost:8083/api/v1/orders"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}