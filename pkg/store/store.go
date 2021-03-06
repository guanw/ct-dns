package store

import (
	"encoding/json"

	"github.com/guanw/ct-dns/pkg/logging"
	storageInterface "github.com/guanw/ct-dns/storage"
	"github.com/pkg/errors"
)

type store struct {
	Client storageInterface.Client
}

// NewStore creates new store instance
func NewStore(client storageInterface.Client) Store {
	logging.GetLogger().Info("Creating new store...")
	return &store{
		Client: client,
	}
}

func (s *store) GetService(serviceName string) ([]string, error) {
	res, err := s.Client.Get(serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "Service Name not found")
	}
	hosts, err := unmarshalStrToHosts(res)
	if err != nil {
		return nil, errors.Wrap(err, "UnmarshalStrToHosts failed")
	}
	return hosts, nil
}

func (s *store) UpdateService(serviceName, operation, host string) error {
	var err error
	if operation == "add" {
		err = s.Client.Create(serviceName, host)
	} else if operation == "delete" {
		err = s.Client.Delete(serviceName, host)
	}
	return err
}

func unmarshalStrToHosts(input string) ([]string, error) {
	var hosts []string
	if err := json.Unmarshal([]byte(input), &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}
