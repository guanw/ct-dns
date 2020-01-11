package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/guanw/ct-dns/storage"
)

// NewFactory creates memory storage client
func NewFactory() (storage.Client, error) {
	pool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
	return NewClient(pool), nil
}
