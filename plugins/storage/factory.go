package storage

import (
	"errors"

	config "github.com/guanw/ct-dns/cmd"

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
	Initialize() (storage.Client, error)
}

type factory struct {
	StorageType string
	V           *viper.Viper
	Cfg         config.Config
}

func (f *factory) Initialize() (storage.Client, error) {
	switch f.StorageType {
	case memoryStorageType:
		return memory.NewFactory()
	case etcdStorageType:
		return etcd.NewFactory(f.V)
	case dynamodbStorageType:
		return dynamodb.NewFactory(f.V)
	case redisStorageType:
		return redis.NewFactory(f.V, f.Cfg)
	default:
		return nil, errors.New("Failed to initialize storage factory")
	}
}

// NewFactory creates storage factory instance
func NewFactory(v *viper.Viper, cfg config.Config) Factory {
	return &factory{
		V:           v,
		Cfg:         cfg,
		StorageType: v.GetString("storage-type"),
	}
}
