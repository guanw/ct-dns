package store

import (
	"github.com/pkg/errors"
)

type retryHandler struct {
	Store             Store
	MaximumRetryTimes int
	Metrics           *Metrics
}

// NewRetryHandler initializes RetryHandler
func NewRetryHandler(maximumRetryTimes int, store Store, metrics *Metrics) Store {
	return &retryHandler{
		MaximumRetryTimes: maximumRetryTimes,
		Store:             store,
		Metrics:           metrics,
	}
}

// GetService fires inner Store maximum times until succeeded
func (r *retryHandler) GetService(serviceName string) ([]string, error) {
	var err error
	var res []string
	for i := 0; i < r.MaximumRetryTimes; i++ {
		if res, err = r.Store.GetService(serviceName); err == nil {
			return res, nil
		}
		r.Metrics.GetServiceRetryAttempts.Inc()
	}
	r.Metrics.GetServiceRetryExhausted.Inc()
	return nil, errors.Wrap(err, "Failed to GetService with RetryHandler")
}

// UpdateService fires inner Store maximum times until succeeded
func (r *retryHandler) UpdateService(serviceName, operation, Host string) error {
	var err error
	for i := 0; i < r.MaximumRetryTimes; i++ {
		if err = r.Store.UpdateService(serviceName, operation, Host); err == nil {
			return nil
		}
		r.Metrics.PostServiceRetryAttempts.Inc()
	}
	r.Metrics.PostServiceRetryExhausted.Inc()
	return errors.Wrap(err, "Failed to PostService with RetryHandler")
}
