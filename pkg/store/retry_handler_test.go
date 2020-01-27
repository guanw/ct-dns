package store

import (
	"errors"
	"testing"

	"github.com/guanw/ct-dns/pkg/store/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRetryHandler_GetService(t *testing.T) {
	mockStore := &mocks.Store{}
	mockStore.On("GetService", "valid-service").Return([]string{"192.0.0.1:8081"}, nil)
	mockStore.On("GetService", "error-service").Return(nil, errors.New("new error"))
	retryHandler := NewRetryHandler(5, mockStore)
	tests := []struct {
		ExpectError bool
		ServiceName string
	}{
		{
			ServiceName: "valid-service",
			ExpectError: false,
		},
		{
			ServiceName: "error-service",
			ExpectError: true,
		},
	}
	for _, test := range tests {
		_, err := retryHandler.GetService(test.ServiceName)
		if test.ExpectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestRetryHandler_PostService(t *testing.T) {
	mockStore := &mocks.Store{}
	mockStore.On("UpdateService", "service", "add", "192.0.0.1:8081").Return(nil)
	mockStore.On("UpdateService", "service", "invalid-operation", "xxx").Return(errors.New("new error"))
	retryHandler := NewRetryHandler(5, mockStore)
	tests := []struct {
		ExpectError bool
		ServiceName string
		Operation   string
		Host        string
	}{
		{
			ServiceName: "service",
			Operation:   "add",
			Host:        "192.0.0.1:8081",
			ExpectError: false,
		},
		{
			ServiceName: "service",
			Operation:   "invalid-operation",
			Host:        "xxx",
			ExpectError: true,
		},
	}
	for _, test := range tests {
		err := retryHandler.UpdateService(test.ServiceName, test.Operation, test.Host)
		if test.ExpectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
