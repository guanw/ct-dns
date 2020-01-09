package storage

// Client defines interface for set/get operation
type Client interface {
	Create(key, value string) error
	Get(key string) (string, error)
	Delete(key, value string) error
}
