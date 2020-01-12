package redis

import (
	"encoding/json"
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/guanw/ct-dns/storage"
	"github.com/pkg/errors"
)

// Pool defines interface for redis.Pool
type Pool interface {
	Get() redis.Conn
}

// Client defines redis client for Create/Get/Delete operations
type Client struct {
	Pool Pool
	lock sync.Mutex
}

// NewClient creates new redis client
func NewClient(pool Pool) storage.Client {
	return &Client{
		Pool: pool,
	}
}

// Create create new key/value pair
func (c *Client) Create(key, value string) error {
	ins := c.Pool.Get()
	defer ins.Close()
	c.lock.Lock()
	_, err := ins.Do("SADD", key, value)
	c.lock.Unlock()
	return err
}

// Get gets hosts under key
func (c *Client) Get(key string) (string, error) {
	ins := c.Pool.Get()
	defer ins.Close()
	c.lock.Lock()
	reply, err := ins.Do("SMEMBERS", key)
	c.lock.Unlock()
	res, err := redis.Strings(reply, err)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get member from key")
	}
	jsonized, _ := json.Marshal(res)
	return string(jsonized), nil
}

// Delete deletes service & host combination
func (c *Client) Delete(key, value string) error {
	ins := c.Pool.Get()
	defer ins.Close()
	c.lock.Lock()
	_, err := ins.Do("SREM", key, value)
	c.lock.Unlock()
	return err
}
