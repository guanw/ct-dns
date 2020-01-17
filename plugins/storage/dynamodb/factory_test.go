package dynamodb

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_NewFactory(t *testing.T) {
	v := viper.New()
	_, err := NewFactory(v)
	assert.NoError(t, err)
}
