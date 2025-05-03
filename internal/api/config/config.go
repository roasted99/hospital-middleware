package config
import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPort string
	DBHost string	
	DBUser string
	DBPassword string
	DBName string
}

func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	envFile := filepath.Join("internal", "api", "config", ".env")
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
	}

	config := &Config{
		DBPort: os.Getenv("DB_PORT"),
		DBHost: os.Getenv("DB_HOST"),
		DBUser: os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
	}

	return config, nil
}

