package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/guanw/ct-dns/pkg/logging"
	"github.com/guanw/ct-dns/storage"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type builder struct {
	Endpoint string
}

func initFromViper(v *viper.Viper) *builder {
	return &builder{
		Endpoint: v.GetString("redis-endpoint"),
	}
}

// NewFactory creates storage client with redis.Pool
func NewFactory(v *viper.Viper) (storage.Client, error) {
	b := initFromViper(v)
	pool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", b.Endpoint)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to create redis pool")
			}
			return c, err
		},
	}
	logging.GetLogger().WithField("Endpoint", b.Endpoint).Info("Creating redis pool")
	return NewClient(pool), nil
}
