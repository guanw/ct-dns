package store

import (
	"github.com/pkg/errors"
)

type retryHandler struct {
	Store             Store
	MaximumRetryTimes int
}

// NewRetryHandler initializes RetryHandler
func NewRetryHandler(maximumRetryTimes int, store Store) Store {
	return &retryHandler{
		MaximumRetryTimes: maximumRetryTimes,
		Store:             store,
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
	}
	return nil, errors.Wrap(err, "Failed to GetService with RetryHandler")
}

// UpdateService fires inner Store maximum times until succeeded
func (r *retryHandler) UpdateService(serviceName, operation, Host string) error {
	var err error
	for i := 0; i < r.MaximumRetryTimes; i++ {
		if err = r.Store.UpdateService(serviceName, operation, Host); err == nil {
			return nil
		}
	}
	return errors.Wrap(err, "Failed to PostService with RetryHandler")
}
