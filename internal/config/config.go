package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Load loads environment variables from .env file
func Load() error {
	return godotenv.Load()
}

// GetDBConfig returns database configuration from environment variables
func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "hospital_middleware"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// GetJWTSecret returns JWT secret from environment variables
func GetJWTSecret() string {
	return getEnv("JWT_SECRET", "6VIY496XKzMoZkj0dJWaMkrh0+oD1pbpIky7nu27QzFsLm0JQOcNllzKRXv8")
}

// GetHospitalAURL returns Hospital A API URL
func GetHospitalAURL() string {
	return getEnv("HOSPITAL_A_URL", "https://hospital-a.api.co.th")
}

// DBConfig represents database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

