package etcd

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/guanw/ct-dns/pkg/storage"
	"go.etcd.io/etcd/client"
)

// Client defines api client for set/get operations
type Client struct {
	API  client.KeysAPI
	lock sync.Mutex
}

// NewClient creates new api client
func NewClient(api client.KeysAPI) storage.Client {
	return &Client{
		API: api,
	}
}

// Create sets new /key/value directory
func (c *Client) Create(key, value string) error {
	c.lock.Lock()
	_, err := c.API.Set(context.Background(), "/"+key+"/"+value, "", &client.SetOptions{
		PrevExist: client.PrevIgnore,
	})
	c.lock.Unlock()
	return err
}

// Get gets values under /key directory
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
func (c *Client) Delete(key, value string) error {
	c.lock.Lock()
	_, err := c.API.Delete(context.Background(), "/"+key+"/"+value, &client.DeleteOptions{
		Recursive: true,
	})
	c.lock.Unlock()
	return err
}
