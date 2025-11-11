package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type HealthChecker struct {
	rabbitmqURL        string
	redis              *redis.Client
	templateServiceURL string
	httpClient         *http.Client
}

func NewHealthChecker(rabbitmqURL string, redis *redis.Client, templateServiceURL string) *HealthChecker {
	return &HealthChecker{
		rabbitmqURL:        rabbitmqURL,
		redis:              redis,
		templateServiceURL: templateServiceURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check RabbitMQ
	conn, err := amqp.Dial(h.rabbitmqURL)
	if err != nil {
		checks["rabbitmq"] = "unhealthy: " + err.Error()
		allHealthy = false
	} else {
		conn.Close()
		checks["rabbitmq"] = "healthy"
	}

	// Check Redis
	if err := h.redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = "unhealthy: " + err.Error()
		allHealthy = false
	} else {
		checks["redis"] = "healthy"
	}

	// Check Template Service
	templateURL := fmt.Sprintf("%s/health", h.templateServiceURL)
	req, _ := http.NewRequestWithContext(ctx, "GET", templateURL, nil)
	resp, err := h.httpClient.Do(req)
	if err != nil {
		checks["template_service"] = "unhealthy: " + err.Error()
		allHealthy = false
	} else {
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			checks["template_service"] = "healthy"
		} else {
			checks["template_service"] = fmt.Sprintf("unhealthy: status %d", resp.StatusCode)
			allHealthy = false
		}
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
