package memory

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/guanw/ct-dns/pkg/storage"
)

type memory interface {
	put(key, value string)
	get(key string) ([]string, error)
	delete(key, value string)
}

type memoryInstance struct {
	data map[string]map[string]interface{}
	lock sync.Mutex
}

func newMemory() memory {
	return &memoryInstance{
		data: make(map[string]map[string]interface{}),
	}
}

func (m *memoryInstance) put(key, value string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, found := m.data[key]
	if !found {
		m.data[key] = make(map[string]interface{})
	}
	m.data[key][value] = true
	return
}

func (m *memoryInstance) get(key string) ([]string, error) {
	m.lock.Lock()
	m.lock.Unlock()
	_, found := m.data[key]
	if !found {
		return nil, errors.New("key not found")
	}
	var res []string
	for nestedK, _ := range m.data[key] {
		res = append(res, nestedK)
	}
	return res, nil
}

func (m *memoryInstance) delete(key, value string) {
	m.lock.Lock()
	m.lock.Unlock()
	nestedMap, found := m.data[key]
	if !found {
		return
	}
	_, found = nestedMap[value]
	if found {
		delete(nestedMap, value)
		if len(nestedMap) == 0 {
			delete(m.data, key)
		}
	}
}

// Client defines storage client using memory
type Client struct {
	m    memory
	lock sync.Mutex
}

// NewClient creates new memory client
func NewClient() storage.Client {
	return &Client{
		m: newMemory(),
	}
}

// Create create new key/value pair
func (c *Client) Create(key, value string) error {
	c.m.put(key, value)
	return nil
}

// Get gets hosts under key
func (c *Client) Get(key string) (string, error) {
	res, err := c.m.get(key)
	if err != nil {
		return "", err
	}
	jsonized, _ := json.Marshal(res)
	return string(jsonized), nil
}

// Delete deletes service & host combination
func (c *Client) Delete(key, value string) error {
	c.m.delete(key, value)
	return nil
}
