package store

import (
	"errors"
	"testing"

	"github.com/guanw/ct-dns/pkg/store/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	metrics      = InitializeMetrics()
	maximumRetry = 10
)

func TestRetryHandler_GetService(t *testing.T) {
	mockStore := &mocks.Store{}
	mockStore.On("GetService", "valid-service").Return([]string{"192.0.0.1:8081"}, nil)
	mockStore.On("GetService", "error-service").Return(nil, errors.New("new error"))
	retryHandler := NewRetryHandler(maximumRetry, mockStore, metrics)
	tests := []struct {
		ExpectError            bool
		ServiceName            string
		ExpectedRetryAttempts  float64
		ExpectedRetryExhausted float64
	}{
		{
			ServiceName:            "valid-service",
			ExpectError:            false,
			ExpectedRetryAttempts:  0.0,
			ExpectedRetryExhausted: 0.0,
		},
		{
			ServiceName:            "error-service",
			ExpectError:            true,
			ExpectedRetryAttempts:  float64(maximumRetry),
			ExpectedRetryExhausted: 1.0,
		},
	}
	for _, test := range tests {
		_, err := retryHandler.GetService(test.ServiceName)
		if test.ExpectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.ExpectedRetryAttempts, testutil.ToFloat64(metrics.GetServiceRetryAttempts))
		assert.Equal(t, test.ExpectedRetryExhausted, testutil.ToFloat64(metrics.GetServiceRetryExhausted))
	}
}

func TestRetryHandler_PostService(t *testing.T) {
	mockStore := &mocks.Store{}
	mockStore.On("UpdateService", "service", "add", "192.0.0.1:8081").Return(nil)
	mockStore.On("UpdateService", "service", "invalid-operation", "xxx").Return(errors.New("new error"))
	retryHandler := NewRetryHandler(maximumRetry, mockStore, metrics)
	tests := []struct {
		ExpectError            bool
		ServiceName            string
		Operation              string
		Host                   string
		ExpectedRetryExhausted float64
		ExpectedRetryAttempts  float64
	}{
		{
			ServiceName:            "service",
			Operation:              "add",
			Host:                   "192.0.0.1:8081",
			ExpectError:            false,
			ExpectedRetryExhausted: 0.0,
			ExpectedRetryAttempts:  0.0,
		},
		{
			ServiceName:            "service",
			Operation:              "invalid-operation",
			Host:                   "xxx",
			ExpectError:            true,
			ExpectedRetryExhausted: 1.0,
			ExpectedRetryAttempts:  float64(maximumRetry),
		},
	}
	for _, test := range tests {
		err := retryHandler.UpdateService(test.ServiceName, test.Operation, test.Host)
		if test.ExpectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.ExpectedRetryAttempts, testutil.ToFloat64(metrics.PostServiceRetryAttempts))
		assert.Equal(t, test.ExpectedRetryExhausted, testutil.ToFloat64(metrics.PostServiceRetryExhausted))
	}
}
