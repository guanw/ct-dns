package storage

import (
	"testing"

	config "github.com/guanw/ct-dns/cmd"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_NewFactory(t *testing.T) {
	tests := []struct {
		expectedErr bool
		factoryType string
	}{
		{
			expectedErr: false,
			factoryType: "dynamodb",
		},
		{
			expectedErr: false,
			factoryType: "etcd",
		},
		{
			expectedErr: false,
			factoryType: "memory",
		},
		{
			expectedErr: false,
			factoryType: "redis",
		},
		{
			expectedErr: true,
			factoryType: "unknown",
		},
	}

	for _, test := range tests {
		v := viper.New()
		flagSet := new(pflag.FlagSet)
		flagSet.String("storage-type", test.factoryType, "--storage-type <name>")
		v.BindPFlags(flagSet)
		f := NewFactory(v, config.Config{})
		_, err := f.Initialize()
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
