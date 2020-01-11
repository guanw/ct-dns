package redis

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/guanw/ct-dns/storage"
	"github.com/pkg/errors"
)

type Client struct {
	Pool *redis.Pool
}

func NewClient(pool *redis.Pool) storage.Client {
	return &Client{
		Pool: pool,
	}
}

// Create create new key/value pair
func (c *Client) Create(key, value string) error {
	ins := c.Pool.Get()
	defer ins.Close()
	_, err := ins.Do("SADD", key, value)
	return err
}

// Get gets hosts under key
func (c *Client) Get(key string) (string, error) {
	ins := c.Pool.Get()
	defer ins.Close()
	reply, err := ins.Do("SMEMBERS", key)
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
	_, err := ins.Do("SREM", key, value)
	return err
}
