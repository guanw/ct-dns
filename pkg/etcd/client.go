package etcd

import (
	"context"

	"go.etcd.io/etcd/client"
)

// APIClient defineds api client for set/get operations
type APIClient struct {
	API client.KeysAPI
}

// NewAPIClient creates new api client
func NewAPIClient(api client.KeysAPI) *APIClient {
	return &APIClient{
		API: api,
	}
}

// CreateOrSet create new key/value pair or set existing key
func (c *APIClient) CreateOrSet(key string, value string) error {
	_, err := c.API.Create(context.Background(), key, value)
	if err != nil {
		// Set key "/foo" to value "bar".
		_, err = c.API.Set(context.Background(), key, value, nil)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// Get gets value with key
func (c *APIClient) Get(key string) (string, error) {
	resp, err := c.API.Get(context.Background(), key, nil)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}
