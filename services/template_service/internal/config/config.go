package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Cache    CacheConfig
}

type ServerConfig struct {
	Port     string
	LogLevel string
	Env      string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type CacheConfig struct {
	TTL int // seconds
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Try to read .env file, but don't fail if it doesn't exist
	// Environment variables take precedence
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only return error if it's not a "file not found" error
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// File not found is OK - we'll use environment variables
	}

	config := &Config{
		Server: ServerConfig{
			Port:     getEnvOrDefault("SERVER_PORT", "8081"),
			LogLevel: getEnvOrDefault("LOG_LEVEL", "info"),
			Env:      getEnvOrDefault("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL: getEnvOrDefault("DATABASE_URL", "postgres://template_user:template_pass@localhost:5432/template_db?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnvOrDefault("REDIS_URL", "redis://localhost:6379"),
		},
		Cache: CacheConfig{
			TTL: 600, // 10 minutes
		},
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
