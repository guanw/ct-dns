package storage

import "testing"

import "github.com/stretchr/testify/assert"

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
			expectedErr: true,
			factoryType: "unknown",
		},
	}
	for _, test := range tests {
		f := NewFactory()
		_, err := f.Initialize(test.factoryType)
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
