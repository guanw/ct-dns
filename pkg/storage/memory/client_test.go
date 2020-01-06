package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InsertNewKey(t *testing.T) {
	m := NewClient()
	m.Create("dummy-service", "192.0.0.1")
	res, err := m.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1"]`, res)

	m.Create("dummy-service", "192.0.0.2")
	res, err = m.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1","192.0.0.2"]`, res)
}
func Test_InsertExistingKeys(t *testing.T) {
	m := NewClient()
	m.Create("dummy-service", "192.0.0.1")
	res, err := m.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1"]`, res)

	m.Create("dummy-service", "192.0.0.1")
	res, err = m.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.1"]`, res)
}

func Test_DeleteOnlyExistingKey(t *testing.T) {
	m := NewClient()
	m.Create("dummy-service", "192.0.0.1")
	m.Delete("dummy-service", "192.0.0.1")
	_, err := m.Get("dummy-service")
	assert.Error(t, err)
}

func Test_DeleteExistingKey(t *testing.T) {
	m := NewClient()
	m.Create("dummy-service", "192.0.0.1")
	m.Create("dummy-service", "192.0.0.2")
	m.Delete("dummy-service", "192.0.0.1")
	res, err := m.Get("dummy-service")
	assert.NoError(t, err)
	assert.Equal(t, `["192.0.0.2"]`, res)
}

func Test_DeleteNonExistingFirstKey(t *testing.T) {
	m := NewClient()
	m.Delete("dummy-service", "192.0.0.1")
	res, err := m.Get("dummy-service")
	assert.Error(t, err)
	assert.Equal(t, ``, res)
}
