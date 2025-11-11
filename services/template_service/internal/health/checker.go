package health

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthChecker struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewHealthChecker(db *pgxpool.Pool, redis *redis.Client) *HealthChecker {
	return &HealthChecker{
		db:    db,
		redis: redis,
	}
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Checks    map[string]string `json:"checks"`
	Timestamp string            `json:"timestamp"`
}

func (h *HealthChecker) Check() HealthStatus {
	checks := make(map[string]string)
	allHealthy := true

	// Check PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		checks["postgres"] = "unhealthy: " + err.Error()
		allHealthy = false
	} else {
		checks["postgres"] = "healthy"
	}

	// Check Redis
	if err := h.redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = "unhealthy: " + err.Error()
		allHealthy = false
	} else {
		checks["redis"] = "healthy"
	}

	status := "healthy"
	if !allHealthy {
		status = "unhealthy"
	}

	return HealthStatus{
		Status:    status,
		Checks:    checks,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
