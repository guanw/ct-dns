package storage

import (
	"errors"

	"github.com/guanw/ct-dns/plugins/storage/dynamodb"
	"github.com/guanw/ct-dns/plugins/storage/etcd"
	"github.com/guanw/ct-dns/plugins/storage/memory"
	"github.com/guanw/ct-dns/storage"
)

const (
	memoryStorageType   = "memory"
	dynamodbStorageType = "dynamodb"
	etcdStorageType     = "etcd"
)

// Factory defines interface for factory
type Factory interface {
	Initialize(factoryType string) (storage.Client, error)
}

type factory struct {
}

func (f *factory) Initialize(factoryType string) (storage.Client, error) {
	switch factoryType {
	case memoryStorageType:
		return memory.NewFactory()
	case etcdStorageType:
		return etcd.NewFactory()
	case dynamodbStorageType:
		return dynamodb.NewFactory()
	default:
		return nil, errors.New("Failed to initialize storage factory")
	}
}

// NewFactory creates storage factory instance
func NewFactory() Factory {
	return &factory{}
}
