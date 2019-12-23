package etcd

import (
	"context"
	"errors"
	"testing"

	"github.com/clicktherapeutics/ct-dns/pkg/etcd"
	"github.com/clicktherapeutics/ct-dns/pkg/etcd/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	client "go.etcd.io/etcd/client"
)

func Test_SetKeyValueNonError(t *testing.T) {
	assert.True(t, true)
	api := &mocks.KeysAPI{}
	api.On("Create", context.Background(), "dummy-service", "[192.0.0.1, 192.0.0.2]").Return(&client.Response{
		Node: &client.Node{
			Key:   "dummy-service",
			Value: "[192.0.0.1, 192.0.0.2]",
		},
	}, nil)
	client := etcd.NewAPIClient(api)
	err := client.CreateOrSet("dummy-service", "[192.0.0.1, 192.0.0.2]")
	assert.NoError(t, err)
}

func Test_SetExistingKeyNoError(t *testing.T) {
	api := &mocks.KeysAPI{}
	api.On("Create", context.Background(), "dummy-service", "[192.0.0.1]").Return(nil, errors.New("dummy-service existed"))
	api.On("Set", context.Background(), "dummy-service", "[192.0.0.1]", mock.Anything).Return(&client.Response{
		Node: &client.Node{
			Key:   "dummy-service",
			Value: "[192.0.0.1]",
		},
	}, nil)
	client := etcd.NewAPIClient(api)
	err := client.CreateOrSet("dummy-service", "[192.0.0.1]")
	assert.NoError(t, err)
}

func Test_Get(t *testing.T) {
	api := &mocks.KeysAPI{}
	api.On("Get", context.Background(), "dummy-service", mock.Anything).Return(&client.Response{
		Node: &client.Node{
			Value: "[192.0.0.1]",
		},
	}, nil)
	client := etcd.NewAPIClient(api)
	res, err := client.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, "[192.0.0.1]", res)

	api.On("Get", context.Background(), "unseen-service", mock.Anything).Return(nil, errors.New("not seen"))
	res, err = client.Get("unseen-service")
	assert.Error(t, err)
	assert.Equal(t, "", res)
}
