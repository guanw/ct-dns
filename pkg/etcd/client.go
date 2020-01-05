package etcd

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"go.etcd.io/etcd/client"
)

// ETCDClient defines interface for set/get operation
type ETCDClient interface {
	Create(key, value string) error
	Get(key string) (string, error)
	Delete(key string, value string) error
}

// Client defineds api client for set/get operations
type Client struct {
	API  client.KeysAPI
	lock sync.Mutex
}

// NewClient creates new api client
func NewClient(api client.KeysAPI) ETCDClient {
	return &Client{
		API: api,
	}
}

// Create sets new key/value pair or set existing key
func (c *Client) Create(key, value string) error {
	c.lock.Lock()
	_, err := c.API.Set(context.Background(), "/"+key+"/"+value, "", &client.SetOptions{
		PrevExist: client.PrevIgnore,
	})
	c.lock.Unlock()
	return err
}

// Get gets value with key
func (c *Client) Get(key string) (string, error) {
	c.lock.Lock()
	resp, err := c.API.Get(context.Background(), "/"+key, &client.GetOptions{
		Recursive: true,
	})
	c.lock.Unlock()
	if err != nil {
		return "", err
	}
	var res []string
	for index := range resp.Node.Nodes {
		res = append(res, strings.Split(resp.Node.Nodes[index].Key, "/")[2])
	}
	json, _ := json.Marshal(res)
	return string(json), nil
}

// Delete delete key recursively
func (c *Client) Delete(key string, value string) error {
	c.lock.Lock()
	_, err := c.API.Delete(context.Background(), "/"+key+"/"+value, &client.DeleteOptions{
		Recursive: true,
	})
	c.lock.Unlock()
	return err
}
