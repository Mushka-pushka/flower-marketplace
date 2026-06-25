package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	RabbitMQURL  string
	ValkeyHost   string
	ValkeyPort   string
	ValkeyDB     int
	LogLevel     string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8082"),
		DBHost:       getEnv("DB_HOST", "127.0.0.1"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "flower_user"),
		DBPassword:   getEnv("DB_PASSWORD", "flower_pass"),
		DBName:       getEnv("DB_NAME", "flower_marketplace"),
		JWTSecret:    getEnv("JWT_SECRET", "your_super_secret_key_here"),
		RabbitMQURL:  getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ValkeyHost:   getEnv("VALKEY_HOST", "localhost"),
		ValkeyPort:   getEnv("VALKEY_PORT", "6379"),
		ValkeyDB:     getEnvAsInt("VALKEY_DB", 0),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}