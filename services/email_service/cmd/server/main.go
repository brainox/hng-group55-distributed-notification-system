package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/circuit"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/config"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/health"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/idempotency"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/queue"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/retry"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/sender"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/template"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.Server.LogLevel); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Log.Info("starting email service", zap.String("version", "1.0.0"))

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.URL,
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Log.Fatal("failed to connect to Redis", zap.Error(err))
	}
	logger.Log.Info("connected to Redis")

	// Initialize components
	idempotencyTTL := time.Duration(24) * time.Hour
	idempotencyChecker := idempotency.NewChecker(redisClient, idempotencyTTL)
	retryHandler := retry.NewHandler(cfg.Retry.MaxAttempts, cfg.Retry.BackoffBase)
	circuitBreaker := circuit.NewBreaker(
		"email-sender",
		1,
		0,
		time.Duration(cfg.CircuitBreaker.Timeout)*time.Second,
	)

	// Initialize email sender
	var emailSender sender.EmailSender
	if cfg.Email.SendGrid.APIKey != "" {
		emailSender, err = sender.NewSendGridSender(cfg.Email.SendGrid)
		if err != nil {
			logger.Log.Fatal("failed to create SendGrid sender", zap.Error(err))
		}
		logger.Log.Info("using SendGrid email sender")
	} else {
		emailSender, err = sender.NewSMTPSender(cfg.Email.SMTP)
		if err != nil {
			logger.Log.Fatal("failed to create SMTP sender", zap.Error(err))
		}
		logger.Log.Info("using SMTP email sender", zap.String("host", cfg.Email.SMTP.Host))
	}

	// Initialize template client
	templateClient := template.NewClient(cfg.TemplateService.URL, redisClient)

	// Initialize queue publisher
	publisher, err := queue.NewPublisher(cfg.RabbitMQ.URL, cfg.RabbitMQ.StatusQueueName)
	if err != nil {
		logger.Log.Fatal("failed to create publisher", zap.Error(err))
	}
	defer publisher.Close()

	// Initialize queue consumer
	consumer, err := queue.NewConsumer(queue.ConsumerConfig{
		URL:            cfg.RabbitMQ.URL,
		QueueName:      cfg.RabbitMQ.QueueName,
		WorkerCount:    cfg.RabbitMQ.WorkerCount,
		TemplateClient: templateClient,
		EmailSender:    emailSender,
		Publisher:      publisher,
		Idempotency:    idempotencyChecker,
		RetryHandler:   retryHandler,
		CircuitBreaker: circuitBreaker,
	})
	if err != nil {
		logger.Log.Fatal("failed to create consumer", zap.Error(err))
	}
	defer consumer.Stop()

	// Start consumer
	if err := consumer.Start(); err != nil {
		logger.Log.Fatal("failed to start consumer", zap.Error(err))
	}

	// Initialize health checker
	healthChecker := health.NewHealthChecker(cfg.RabbitMQ.URL, redisClient, cfg.TemplateService.URL)

	// Setup HTTP server for health checks
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		status := healthChecker.Check()
		httpStatus := http.StatusOK
		if status.Status == "unhealthy" {
			httpStatus = http.StatusServiceUnavailable
		}
		c.JSON(httpStatus, status)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start HTTP server
	go func() {
		logger.Log.Info("starting HTTP server", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down gracefully...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("server stopped")
}
