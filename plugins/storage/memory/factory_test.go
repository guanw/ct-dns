package memory

import "testing"

import "github.com/stretchr/testify/assert"

func Test_NewFactory(t *testing.T) {
	_, err := NewFactory()
	assert.NoError(t, err)
}
