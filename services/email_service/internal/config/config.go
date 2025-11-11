package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server          ServerConfig
	RabbitMQ        RabbitMQConfig
	Redis           RedisConfig
	TemplateService TemplateServiceConfig
	Email           EmailConfig
	Retry           RetryConfig
	CircuitBreaker  CircuitBreakerConfig
}

type ServerConfig struct {
	Port     string
	LogLevel string
	Env      string
}

type RabbitMQConfig struct {
	URL             string
	QueueName       string
	StatusQueueName string
	WorkerCount     int
}

type RedisConfig struct {
	URL string
}

type TemplateServiceConfig struct {
	URL string
}

type EmailConfig struct {
	Provider string // "smtp" or "sendgrid"
	SMTP     SMTPConfig
	SendGrid SendGridConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SendGridConfig struct {
	APIKey string
}

type RetryConfig struct {
	MaxAttempts int
	BackoffBase int // seconds
}

type CircuitBreakerConfig struct {
	Threshold int
	Timeout   int // seconds
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	smtpPort, _ := strconv.Atoi(getEnvOrDefault("SMTP_PORT", "587"))
	workerCount, _ := strconv.Atoi(getEnvOrDefault("WORKER_COUNT", "10"))
	maxRetry, _ := strconv.Atoi(getEnvOrDefault("MAX_RETRY_ATTEMPTS", "5"))
	backoff, _ := strconv.Atoi(getEnvOrDefault("RETRY_BACKOFF_BASE", "1"))
	cbThreshold, _ := strconv.Atoi(getEnvOrDefault("CIRCUIT_BREAKER_THRESHOLD", "5"))
	cbTimeout, _ := strconv.Atoi(getEnvOrDefault("CIRCUIT_BREAKER_TIMEOUT", "30"))

	config := &Config{
		Server: ServerConfig{
			Port:     getEnvOrDefault("SERVER_PORT", "8082"),
			LogLevel: getEnvOrDefault("LOG_LEVEL", "info"),
			Env:      getEnvOrDefault("ENV", "development"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:             getEnvOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			QueueName:       getEnvOrDefault("QUEUE_NAME", "email.queue"),
			StatusQueueName: getEnvOrDefault("STATUS_QUEUE_NAME", "notification.status.queue"),
			WorkerCount:     workerCount,
		},
		Redis: RedisConfig{
			URL: getEnvOrDefault("REDIS_URL", "redis://localhost:6379"),
		},
		TemplateService: TemplateServiceConfig{
			URL: getEnvOrDefault("TEMPLATE_SERVICE_URL", "http://localhost:8081"),
		},
		Email: EmailConfig{
			Provider: getEnvOrDefault("EMAIL_PROVIDER", "smtp"),
			SMTP: SMTPConfig{
				Host:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
				Port:     smtpPort,
				Username: getEnvOrDefault("SMTP_USERNAME", ""),
				Password: getEnvOrDefault("SMTP_PASSWORD", ""),
			},
			SendGrid: SendGridConfig{
				APIKey: getEnvOrDefault("SENDGRID_API_KEY", ""),
			},
		},
		Retry: RetryConfig{
			MaxAttempts: maxRetry,
			BackoffBase: backoff,
		},
		CircuitBreaker: CircuitBreakerConfig{
			Threshold: cbThreshold,
			Timeout:   cbTimeout,
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
