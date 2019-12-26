package etcd

import (
	"context"

	"go.etcd.io/etcd/client"
)

// ETCDClient defines interface for set/get operation
type ETCDClient interface {
	CreateOrSet(key, value string) error
	Get(key string) (string, error)
}

// Client defineds api client for set/get operations
type Client struct {
	API client.KeysAPI
}

// NewClient creates new api client
func NewClient(api client.KeysAPI) ETCDClient {
	return &Client{
		API: api,
	}
}

// CreateOrSet create new key/value pair or set existing key
func (c *Client) CreateOrSet(key, value string) error {
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
func (c *Client) Get(key string) (string, error) {
	resp, err := c.API.Get(context.Background(), key, nil)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}
