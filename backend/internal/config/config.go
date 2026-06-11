package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env           string
	Port          string
	AllowedOrigin string
	GroqAPIKey    string

	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// Auth
	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshExpiry time.Duration
}

// IsProd returns true if the application is running in production mode.
func (c *Config) IsProd() bool {
	return c.Env == "production"
}

// Load reads configuration from environment variables with default fallbacks.
func Load() *Config {
	jwtExpiryRaw := getEnv("JWT_EXPIRY_MINUTES", "15")
	jwtMinutes, err := strconv.Atoi(jwtExpiryRaw)
	if err != nil {
		log.Printf("Warning: invalid JWT_EXPIRY_MINUTES '%s'. Using default: 15", jwtExpiryRaw)
		jwtMinutes = 15
	}
	jwtExpiry := time.Duration(jwtMinutes) * time.Minute

	refreshRaw := getEnv("REFRESH_EXPIRY_DAYS", "30")
	refreshDays, err := strconv.Atoi(refreshRaw)
	if err != nil {
		log.Printf("Warning: invalid REFRESH_EXPIRY_DAYS '%s' (must be a number). Using default: 30", refreshRaw)
		refreshDays = 30
	}
	refreshExpiry := time.Duration(refreshDays) * 24 * time.Hour

	return &Config{
		Env:           getEnv("ENV", "development"),
		Port:          getEnv("PORT", "8080"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
		GroqAPIKey:    getEnv("GROQ_API_KEY", ""),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "todo_db"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key"),
		JWTExpiry:     jwtExpiry,
		RefreshExpiry: refreshExpiry,
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
