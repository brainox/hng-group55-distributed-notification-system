package circuit

import (
	"time"

	"github.com/sony/gobreaker"
)

func NewBreaker(name string, maxRequests uint32, interval time.Duration, timeout time.Duration) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: maxRequests,
		Interval:    interval,
		Timeout:     timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}
