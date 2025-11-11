package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/config"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/handler"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/health"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/middleware"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/repository"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/service"
	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Server.LogLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Log.Info("starting template service", zap.String("port", cfg.Server.Port))

	// Run database migrations
	if err := runMigrations(cfg.Database.URL); err != nil {
		logger.Log.Fatal("failed to run migrations", zap.Error(err))
	}

	// Connect to PostgreSQL
	dbPool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	logger.Log.Info("connected to PostgreSQL")

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.URL[8:], // Remove "redis://" prefix
	})
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Log.Fatal("failed to connect to Redis", zap.Error(err))
	}

	logger.Log.Info("connected to Redis")

	// Initialize repositories
	templateRepo := repository.NewTemplateRepository(dbPool)
	versionRepo := repository.NewVersionRepository(dbPool)

	// Initialize services
	cacheService := service.NewCacheService(redisClient)
	templateService := service.NewTemplateService(
		templateRepo,
		versionRepo,
		cacheService,
		time.Duration(cfg.Cache.TTL)*time.Second,
	)

	// Initialize handlers
	templateHandler := handler.NewTemplateHandler(templateService)

	// Initialize health checker
	healthChecker := health.NewHealthChecker(dbPool, redisClient)

	// Setup Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		status := healthChecker.Check()
		statusCode := http.StatusOK
		if status.Status != "healthy" {
			statusCode = http.StatusServiceUnavailable
		}
		c.JSON(statusCode, status)
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		templates := v1.Group("/templates")
		{
			templates.POST("", templateHandler.CreateTemplate)
			templates.GET("", templateHandler.ListTemplates)
			templates.GET("/:id", templateHandler.GetTemplateByID)
			templates.GET("/key/:key", templateHandler.GetTemplateByKey)
			templates.PUT("/:id", templateHandler.UpdateTemplate)
			templates.DELETE("/:id", templateHandler.DeleteTemplate)

			templates.POST("/:id/versions", templateHandler.CreateVersion)
			templates.GET("/:id/versions", templateHandler.GetVersionHistory)
			templates.POST("/:id/versions/:version_id/publish", templateHandler.PublishVersion)
			templates.POST("/:id/preview", templateHandler.PreviewTemplate)

			templates.POST("/validate", templateHandler.ValidateTemplate)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	logger.Log.Info("template service started", zap.String("port", cfg.Server.Port))

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("server exited")
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Log.Info("database migrations completed")
	return nil
}
