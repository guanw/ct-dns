package etcd

import (
	"context"
	"errors"
	"testing"

	"github.com/guanw/ct-dns/plugins/storage/etcd/mocks"
	"github.com/stretchr/testify/assert"
	client "go.etcd.io/etcd/client"
)

func Test_SetKeyValueNonError(t *testing.T) {
	assert.True(t, true)
	api := &mocks.KeysAPI{}
	api.On("Set", context.Background(), "/dummy-service/[192.0.0.1, 192.0.0.2]", "", &client.SetOptions{
		PrevExist: client.PrevIgnore,
	}).Return(&client.Response{
		Node: &client.Node{
			Key:   "dummy-service",
			Value: "[192.0.0.1, 192.0.0.2]",
		},
	}, nil)
	client := NewClient(api)
	err := client.Create("dummy-service", "[192.0.0.1, 192.0.0.2]")
	assert.NoError(t, err)
}

func Test_Get(t *testing.T) {
	api := &mocks.KeysAPI{}
	api.On("Get", context.Background(), "/dummy-service", &client.GetOptions{
		Recursive: true,
	}).Return(&client.Response{
		Node: &client.Node{
			Nodes: client.Nodes{
				{
					Key: "/dummy-service/192.0.0.1",
				},
				{
					Key: "/dummy-service/192.0.0.2",
				},
			},
		},
	}, nil)
	cli := NewClient(api)
	res, err := cli.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, "[\"192.0.0.1\",\"192.0.0.2\"]", res)

	api.On("Get", context.Background(), "/unseen-service", &client.GetOptions{
		Recursive: true,
	}).Return(nil, errors.New("not seen"))
	res, err = cli.Get("unseen-service")
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func Test_Delete(t *testing.T) {
	api := &mocks.KeysAPI{}
	api.On("Delete", context.Background(), "/dummy-service/192.0.0.1", &client.DeleteOptions{
		Recursive: true,
	}).Return(&client.Response{}, nil)
	cli := NewClient(api)
	err := cli.Delete("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
}
