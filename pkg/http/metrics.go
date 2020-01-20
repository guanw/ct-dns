package http

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics defines all metrics for http server
type Metrics struct {
	GetServiceSuccess prometheus.Counter
	GetServiceFailure prometheus.Counter

	PostServiceSuccess prometheus.Counter
	PostServiceFailure prometheus.Counter

	HealthcheckSuccess prometheus.Counter

	V1RegistrationSuccess prometheus.Counter
	V1RegistrationFailure prometheus.Counter

	V2DiscoverySuccess prometheus.Counter
	V2DiscoveryFailure prometheus.Counter
}

// InitializeMetrics initialize http metrics
func InitializeMetrics() *Metrics {
	return &Metrics{
		GetServiceSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_handler_get_request_success",
		}),
		GetServiceFailure: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_handler_get_request_failure",
		}),
		PostServiceSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_handler_post_request_success",
		}),
		PostServiceFailure: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_handler_post_request_failure",
		}),
		HealthcheckSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_handler_health_check_success",
		}),
		V1RegistrationSuccess: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_handler_v1_registration_success",
		}),
		V1RegistrationFailure: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_handler_v1_registration_failure",
		}),
		V2DiscoverySuccess: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_handler_v2_discovery_success",
		}),
		V2DiscoveryFailure: promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_handler_v2_discovery_failure",
		}),
	}
}
