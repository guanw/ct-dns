package store

import (
	"encoding/json"
	"fmt"

	"github.com/clicktherapeutics/ct-dns/pkg/etcd"
	"github.com/pkg/errors"
)

type Store interface {
	GetService(serviceName string) ([]string, error)
	UpdateService(serviceName, operation, Host string) error
}

type store struct {
	Client etcd.ETCDClient
}

func NewStore(client etcd.ETCDClient) *store {
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
	res, err := s.Client.Get(serviceName)
	if err != nil {
		if operation == "add" {
			err := s.Client.CreateOrSet(serviceName, fmt.Sprintf(`["%s"]`, host))
			if err != nil {
				return errors.Wrap(err, "Failed to create new service entry")
			}
			return nil
		}
		return errors.Wrap(err, "Failed to delete service: Service name not found")
	}
	hosts, err := unmarshalStrToHosts(res)
	if err != nil {
		return errors.Wrap(err, "UnmarshalStrToHosts failed")
	}
	if operation == "add" {
		include, _ := contains(hosts, host)
		if !include {
			hosts = append(hosts, host)
			str, err := marshalHostsToStr(hosts)
			if err != nil {
				return errors.Wrap(err, "Failed to marshal hosts to string")
			}
			err = s.Client.CreateOrSet(serviceName, str)
			if err != nil {
				return errors.Wrap(err, "Failed to add host")
			}
			return nil
		}
		return errors.New("Failed to add Host: Host already existed")
	} else if operation == "delete" {
		include, i := contains(hosts, host)
		if include {
			hosts = append(hosts[:i], hosts[i+1:]...)
			str, err := marshalHostsToStr(hosts)
			if err != nil {
				return errors.Wrap(err, "Failed to marshal hosts to string")
			}
			err = s.Client.CreateOrSet(serviceName, str)
			if err != nil {
				return errors.Wrap(err, "Failed to delete host")
			}
			return nil
		}
		return errors.New("Failed to delete Host: Host not found")
	}
	return nil
}

func unmarshalStrToHosts(input string) ([]string, error) {
	var hosts []string
	if err := json.Unmarshal([]byte(input), &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

func marshalHostsToStr(input []string) (string, error) {
	res, err := json.Marshal(input)
	return string(res), err
}

func contains(s []string, e string) (bool, int) {
	for i, a := range s {
		if a == e {
			return true, i
		}
	}
	return false, -1
}
