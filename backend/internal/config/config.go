package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPassword    string
	DBSSLMode     string
	JWTSecret     string
	JWTExpiry     time.Duration
	Port          string
	AllowedOrigin string
	GroqAPIKey    string
}

func Load() *Config {
	expiryRaw := getEnv("JWT_EXPIRY_HOURS", "24")

	hours, err := strconv.Atoi(expiryRaw)
	if err != nil {
		log.Printf("Warning: invalid JWT_EXPIRY_HOURS '%s' (must be a number). Using default: 24", expiryRaw)
		hours = 24
	}

	return &Config{
		DBHost:        getEnv("DB_HOST", "127.0.0.1"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBName:        getEnv("DB_NAME", "todo_db"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key"),
		JWTExpiry:     time.Duration(hours) * time.Hour,
		Port:          getEnv("PORT", "8080"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
		GroqAPIKey:    getEnv("GROQ_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
