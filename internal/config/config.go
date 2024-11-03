package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// Config represents the application configuration structure
type Config struct {
	DSN          string
	Port         string
	JWTSecretKey string
	Environment  string
}

func Load() (*Config, error) {
	config := &Config{
		DSN:          os.Getenv("DSN"),
		Port:         os.Getenv("PORT"),
		JWTSecretKey: os.Getenv("JWY_SECRET_KEY"),
		Environment:  os.Getenv("ENVIRONMENT"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validateConfig checks for required configuration fields and returns an error if any are missing
func (c *Config) validate() error {
	if c.DSN == "" {
		return fmt.Errorf("DSN is required")
	}
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if c.JWTSecretKey == "" {
		return fmt.Errorf("JWT_SECRET_KEY is required")
	}
	if c.Environment == "" {
		return fmt.Errorf("ENVIRONMENT is required")
	}
	return nil
}
