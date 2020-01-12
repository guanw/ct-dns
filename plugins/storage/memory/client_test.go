package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InsertNewKey(t *testing.T) {
	m := NewClient()
	err := m.Create("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
	res, err := m.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1"]`, res)

	err = m.Create("dummy-service", "192.0.0.2")
	assert.NoError(t, err)
	res, err = m.Get("dummy-service")
	assert.NoError(t, err)
	assert.True(t, `["192.0.0.1","192.0.0.2"]` == res || `["192.0.0.2","192.0.0.1"]` == res)
}
func Test_InsertExistingKeys(t *testing.T) {
	m := NewClient()
	err := m.Create("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
	res, err := m.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1"]`, res)

	err = m.Create("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
	res, err = m.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1"]`, res)
}

func Test_DeleteOnlyExistingKey(t *testing.T) {
	m := NewClient()
	err := m.Create("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
	err = m.Delete("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
	_, err = m.Get("dummy-service")
	assert.Error(t, err)
}

func Test_DeleteExistingKey(t *testing.T) {
	m := NewClient()
	err := m.Create("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
	err = m.Create("dummy-service", "192.0.0.2")
	assert.NoError(t, err)
	err = m.Delete("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
	res, err := m.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.2"]`, res)
}

func Test_DeleteNonExistingFirstKey(t *testing.T) {
	m := NewClient()
	err := m.Delete("dummy-service", "192.0.0.1")
	assert.NoError(t, err)
	res, err := m.Get("dummy-service")
	assert.Error(t, err)
	assert.Equal(t, ``, res)
}
