package storage

import (
	"errors"

	"github.com/guanw/ct-dns/plugins/storage/dynamodb"
	"github.com/guanw/ct-dns/plugins/storage/etcd"
	"github.com/guanw/ct-dns/plugins/storage/memory"
	"github.com/guanw/ct-dns/plugins/storage/redis"
	"github.com/guanw/ct-dns/storage"
	"github.com/spf13/viper"
)

const (
	memoryStorageType   = "memory"
	dynamodbStorageType = "dynamodb"
	etcdStorageType     = "etcd"
	redisStorageType    = "redis"
)

// Factory defines interface for factory
type Factory interface {
	Initialize(v *viper.Viper, factoryType string) (storage.Client, error)
}

type factory struct {
}

func (f *factory) Initialize(v *viper.Viper, factoryType string) (storage.Client, error) {
	switch factoryType {
	case memoryStorageType:
		return memory.NewFactory()
	case etcdStorageType:
		return etcd.NewFactory()
	case dynamodbStorageType:
		return dynamodb.NewFactory(v)
	case redisStorageType:
		return redis.NewFactory()
	default:
		return nil, errors.New("Failed to initialize storage factory")
	}
}

// NewFactory creates storage factory instance
func NewFactory() Factory {
	return &factory{}
}
