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

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid SERVER_PORT: %v", err)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "birthday_db"),
		ServerPort: serverPort,
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