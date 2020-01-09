package memory

import (
	"github.com/guanw/ct-dns/storage"
)

// NewFactory creates memory storage client
func NewFactory() (storage.Client, error) {
	return NewClient(), nil
}
