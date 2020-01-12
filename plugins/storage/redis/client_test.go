package redis

import (
	"testing"

	"github.com/guanw/ct-dns/plugins/storage/redis/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_SetKeyValueNonError(t *testing.T) {
	p := &mocks.Pool{}
	c := &mocks.Conn{}
	p.On("Get").Return(c)
	c.On("Do", "SADD", "dummy-service", "192.0.0.1").Return(nil, nil)
	c.On("Close").Return(nil)
	client := NewClient(p)
	err := client.Create("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
}

func Test_Get(t *testing.T) {
	p := &mocks.Pool{}
	c := &mocks.Conn{}
	p.On("Get").Return(c)
	c.On("Do", "SMEMBERS", "dummy-service").Return([]interface{}{"192.0.0.1", "192.0.0.2"}, nil)
	c.On("Close").Return(nil)
	client := NewClient(p)
	res, err := client.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1","192.0.0.2"]`, res)
}

func Test_Delete(t *testing.T) {
	p := &mocks.Pool{}
	c := &mocks.Conn{}
	p.On("Get").Return(c)
	c.On("Do", "SREM", "dummy-service", "192.0.0.1").Return(nil, nil)
	c.On("Close").Return(nil)
	client := NewClient(p)
	err := client.Delete("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
}
