package store

import (
	"errors"
	"testing"

	"github.com/guanw/ct-dns/storage/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_GetService(t *testing.T) {
	mockClient := &mocks.Client{}
	mockClient.On("Get", "dummy-service").Return(`["192.0.0.1"]`, nil)
	mockClient.On("Get", "non-exist-service").Return("", nil)
	mockClient.On("Get", "error-service").Return("", errors.New("not found"))
	store := NewStore(mockClient)

	tests := []struct {
		expectedErr      bool
		expectedResponse []string
		serviceName      string
	}{
		{
			serviceName:      "dummy-service",
			expectedErr:      false,
			expectedResponse: []string{"192.0.0.1"},
		},
		{
			serviceName: "non-exist-service",
			expectedErr: true,
		},
		{
			serviceName: "error-service",
			expectedErr: true,
		},
	}

	for _, test := range tests {
		hosts, err := store.GetService(test.serviceName)
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResponse, hosts)
		}
	}
}

func Test_ServiceAddNewHost(t *testing.T) {
	mockClient := &mocks.Client{}
	mockClient.On("Create", "dummy-service", "192.0.0.1").Return(nil)
	store := NewStore(mockClient)

	err := store.UpdateService("dummy-service", "add", "192.0.0.1")
	assert.NoError(t, err)
}

func Test_ServiceDeleteHost(t *testing.T) {
	mockClient := &mocks.Client{}
	mockClient.On("Delete", "dummy-service", "192.0.0.1").Return(nil)
	store := NewStore(mockClient)

	err := store.UpdateService("dummy-service", "delete", "192.0.0.1")
	assert.NoError(t, err)
}

func Test_unmarshalStrToHosts(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr bool
		description string
		expected    []string
	}{
		{
			input:       `["192.0.0.1","192.0.0.2"]`,
			expectedErr: false,
			expected:    []string{"192.0.0.1", "192.0.0.2"},
			description: `input: "[192.0.0.1, 192.0.0.2]"`,
		},
		{
			input:       `""`,
			expectedErr: true,
			description: "input invalid",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			out, err := unmarshalStrToHosts(test.input)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, out)
			}
		})
	}
}
