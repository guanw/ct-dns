package etcd

import (
	"github.com/guanw/ct-dns/storage"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/client"
)

// NewFactory creates new etcd factory
func NewFactory() (storage.Client, error) {
	var endpoints []string
	etcd := []struct {
		Host string
		Port string
	}{
		{
			Host: "10.110.251.205",
			Port: "2379",
		},
	}
	for _, v := range etcd {
		endpoints = append(endpoints, "http://"+v.Host+":"+v.Port)
	}
	etcdCfg := client.Config{
		Endpoints: endpoints,
	}

	c, err := client.New(etcdCfg)
	if err != nil {
		// handle error
		return nil, errors.Wrap(err, "Cannot initialize the etcd client")
	}
	return NewClient(client.NewKeysAPI(c)), nil
}
