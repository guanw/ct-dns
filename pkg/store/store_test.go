package store

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockETCDClient struct {
	m map[string]string
}

func newMockETCDClient() *mockETCDClient {
	return &mockETCDClient{
		m: make(map[string]string),
	}
}

func (c *mockETCDClient) CreateOrSet(key, value string) error {
	c.m[key] = value
	return nil
}

func (c *mockETCDClient) Get(key string) (string, error) {
	val, found := c.m[key]
	if !found {
		return "", errors.New("invalid service")
	}
	return val, nil
}

func Test_GetService(t *testing.T) {
	mockClient := newMockETCDClient()
	mockClient.CreateOrSet("dummy-service", `["192.0.0.1"]`)
	mockClient.CreateOrSet("empty-service", `[]`)
	mockClient.CreateOrSet("invalid-service", ``)
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
			serviceName: "invalid-service",
			expectedErr: true,
		},
		{
			serviceName:      "empty-service",
			expectedErr:      false,
			expectedResponse: []string{},
		},
		{
			serviceName: "non-exist",
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

func Test_ExistingServiceAddNewHost(t *testing.T) {
	mockClient := newMockETCDClient()
	store := NewStore(mockClient)

	err := store.UpdateService("dummy-service", "add", "192.0.0.1")
	assert.NoError(t, err)

	err = store.UpdateService("dummy-service", "add", "192.0.0.2")
	assert.NoError(t, err)

	res, err := store.GetService("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, []string{"192.0.0.1", "192.0.0.2"}, res)
}

func Test_ExistingServiceAddExistingHost(t *testing.T) {
	mockClient := newMockETCDClient()
	store := NewStore(mockClient)

	err := store.UpdateService("dummy-service", "add", "192.0.0.1")
	assert.NoError(t, err)

	err = store.UpdateService("dummy-service", "add", "192.0.0.1")
	assert.True(t, strings.Contains(err.Error(), "Failed to add Host: Host already existed"))

	res, err := store.GetService("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, []string{"192.0.0.1"}, res)
}

func Test_NonExistingServiceDeleteHost(t *testing.T) {
	mockClient := newMockETCDClient()
	store := NewStore(mockClient)

	err := store.UpdateService("non-exist-service", "delete", "192.0.0.1")
	assert.True(t, strings.Contains(err.Error(), "Failed to delete service: Service name not found"))
}

func Test_ExistingServiceDeleteNonExistingHost(t *testing.T) {
	mockClient := newMockETCDClient()
	store := NewStore(mockClient)

	err := store.UpdateService("dummy-service", "add", "192.0.0.1")
	assert.NoError(t, err)

	err = store.UpdateService("dummy-service", "delete", "192.0.0.2")
	assert.True(t, strings.Contains(err.Error(), "Failed to delete Host: Host not found"))

	res, err := store.GetService("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, []string{"192.0.0.1"}, res)
}

func Test_ExistingServiceDeleteExistingHost(t *testing.T) {
	mockClient := newMockETCDClient()
	store := NewStore(mockClient)

	err := store.UpdateService("dummy-service", "add", "192.0.0.1")
	assert.NoError(t, err)

	err = store.UpdateService("dummy-service", "add", "192.0.0.2")
	assert.NoError(t, err)

	err = store.UpdateService("dummy-service", "delete", "192.0.0.1")
	assert.NoError(t, err)

	res, err := store.GetService("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, []string{"192.0.0.2"}, res)
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

func Test_marshalHostsToStr(t *testing.T) {
	tests := []struct {
		input       []string
		expectedErr bool
		description string
		expected    string
	}{
		{
			input:       []string{"192.0.0.1", "192.0.0.2"},
			expectedErr: false,
			expected:    `["192.0.0.1","192.0.0.2"]`,
			description: "input: [192.0.0.1, 192.0.0.2]",
		},
		{
			input:       nil,
			expectedErr: false,
			expected:    "null",
			description: "input invalid",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			out, err := marshalHostsToStr(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, out)
		})
	}
}
