package retry

import (
	"math"
	"time"

	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	maxAttempts int
	baseBackoff int // seconds
}

func NewHandler(maxAttempts, baseBackoff int) *Handler {
	return &Handler{
		maxAttempts: maxAttempts,
		baseBackoff: baseBackoff,
	}
}

// ShouldRetry determines if an error is retryable
func (h *Handler) ShouldRetry(err error, attempt int) bool {
	if attempt >= h.maxAttempts {
		return false
	}

	// Check if error is retryable
	errMsg := err.Error()

	// Don't retry on these permanent errors
	permanentErrors := []string{
		"invalid email",
		"template not found",
		"authentication failed",
		"unauthorized",
		"forbidden",
	}

	for _, permErr := range permanentErrors {
		if contains(errMsg, permErr) {
			return false
		}
	}

	// Retry on temporary errors
	return true
}

// CalculateBackoff calculates exponential backoff duration
func (h *Handler) CalculateBackoff(attempt int) time.Duration {
	// Exponential backoff: baseBackoff * 2^attempt
	// Max backoff: 16 seconds
	backoff := float64(h.baseBackoff) * math.Pow(2, float64(attempt))
	if backoff > 16 {
		backoff = 16
	}
	return time.Duration(backoff) * time.Second
}

// Wait waits for the backoff duration
func (h *Handler) Wait(attempt int, correlationID string) {
	backoff := h.CalculateBackoff(attempt)
	logger.Log.Info("retry backoff",
		zap.Int("attempt", attempt),
		zap.Duration("backoff", backoff),
		zap.String("correlation_id", correlationID),
	)
	time.Sleep(backoff)
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr ||
		(len(str) > len(substr) && containsSubstring(str, substr)))
}

func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
