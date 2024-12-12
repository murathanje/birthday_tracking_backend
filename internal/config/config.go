package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort int
	GinMode    string
	APIKey     string
	JWTSecret  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

    return &Config{
        DBHost:     getEnv("DATABASE_HOST", "localhost"),
        DBUser:     getEnv("DATABASE_USER", "postgres"),
        DBPassword: getEnv("DATABASE_PASSWORD", ""),
        DBName:     getEnv("DATABASE_NAME", "birthday_db"),
        ServerPort: getEnvAsInt("SERVER_PORT", 5050),
        GinMode:    getEnv("GIN_MODE", "debug"),
        APIKey:     getEnv("API_KEY", "default-api-key"),
        JWTSecret:  getEnv("JWT_SECRET", "default-jwt-secret"),
    }
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
} 