package config

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
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8081"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "flower_user"),
		DBPassword:  getEnv("DB_PASSWORD", "flower_pass"),
		DBName:      getEnv("DB_NAME", "flower_marketplace"),
		JWTSecret:   getEnv("JWT_SECRET", "your_secret_key_here"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ValkeyHost:  getEnv("VALKEY_HOST", "localhost"),
		ValkeyPort:  getEnv("VALKEY_PORT", "6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
