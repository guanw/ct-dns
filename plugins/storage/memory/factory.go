package memory

import (
	"github.com/guanw/ct-dns/pkg/logging"
	"github.com/guanw/ct-dns/storage"
)

// NewFactory creates memory storage client
func NewFactory() (storage.Client, error) {
	logging.GetLogger().Info("Building memory storage")
	return NewClient(), nil
}
