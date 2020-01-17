package storage

import (
	"testing"

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
		f := NewFactory()
		v := viper.New()
		_, err := f.Initialize(v, test.factoryType)
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
