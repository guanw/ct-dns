package etcd

import (
	"strings"

	"github.com/guanw/ct-dns/pkg/logging"
	"github.com/guanw/ct-dns/storage"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/client"
)

type builder struct {
	Endpoints []string
}

func initFromViper(v *viper.Viper) *builder {
	endpoints := strings.Split(v.GetString("etcd-endpoints"), ",")
	return &builder{
		Endpoints: endpoints,
	}
}

// NewFactory creates new etcd factory
func NewFactory(v *viper.Viper) (storage.Client, error) {
	b := initFromViper(v)
	etcdCfg := client.Config{
		Endpoints: b.Endpoints,
	}

	c, err := client.New(etcdCfg)
	if err != nil {
		// handle error
		return nil, errors.Wrap(err, "Cannot initialize the etcd client")
	}
	logging.GetLogger().WithField("Endpoints", b.Endpoints).Info("Creating etcd session")
	return NewClient(client.NewKeysAPI(c)), nil
}
