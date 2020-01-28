package redis

import (
	"testing"

	config "github.com/guanw/ct-dns/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_NewFactory(t *testing.T) {
	v := viper.New()
	_, err := NewFactory(v, config.Config{})
	assert.NoError(t, err)
}

func Test_initFromViperWithCfg(t *testing.T) {
	expectedHost := "redis-test"
	expectedPort := "6379"
	b := initFromViper(viper.New(), config.Config{
		Redis: config.RedisConfig{
			Host: expectedHost,
			Port: expectedPort,
		},
	})
	assert.Equal(t, expectedHost+":"+expectedPort, b.Endpoint)
}
