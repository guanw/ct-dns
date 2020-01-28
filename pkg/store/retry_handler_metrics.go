package store

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics defines all metrics for retry handler
type Metrics struct {
	GetServiceRetryAttempts  prometheus.Counter
	GetServiceRetryExhausted prometheus.Counter

	PostServiceRetryAttempts  prometheus.Counter
	PostServiceRetryExhausted prometheus.Counter
}

// InitializeMetrics initialize retry metrics
func InitializeMetrics() *Metrics {
	return &Metrics{
		GetServiceRetryAttempts: promauto.NewCounter(prometheus.CounterOpts{
			Name: "get_service_retry_attempts",
		}),
		GetServiceRetryExhausted: promauto.NewCounter(prometheus.CounterOpts{
			Name: "get_service_retry_exhausted",
		}),
		PostServiceRetryAttempts: promauto.NewCounter(prometheus.CounterOpts{
			Name: "update_service_retry_attempts",
		}),
		PostServiceRetryExhausted: promauto.NewCounter(prometheus.CounterOpts{
			Name: "update_service_retry_exhausted",
		}),
	}
}
