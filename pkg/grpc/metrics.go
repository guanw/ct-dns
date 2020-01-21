package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics defines all metrics for grpc server
type Metrics struct {
	GetServiceSuccess prometheus.Counter
	GetServiceFailure prometheus.Counter

	PostServiceSuccess prometheus.Counter
	PostServiceFailure prometheus.Counter
}

// InitializeMetrics initialize grpc metrics
func InitializeMetrics() *Metrics {
	return &Metrics{
		GetServiceSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Name: "grpc_handler_get_service_success",
		}),
		GetServiceFailure: promauto.NewCounter(prometheus.CounterOpts{
			Name: "grpc_handler_get_service_failure",
		}),
		PostServiceSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Name: "grpc_handler_post_service_success",
		}),
		PostServiceFailure: promauto.NewCounter(prometheus.CounterOpts{
			Name: "grpc_handler_post_service_failure",
		}),
	}
}
