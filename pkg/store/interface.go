package store

// Store defines interface
type Store interface {
	GetService(serviceName string) ([]string, error)
	UpdateService(serviceName, operation, Host string) error
}
