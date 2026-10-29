package outgoing

import (
	"customer-management/adapter/outgoing/memory"
	"customer-management/adapter/outgoing/postgres"
	"customer-management/port"
)

// Config holds what the outgoing layer needs to choose a storage.
type Config struct {
	DatabaseURL string
}

// NewCustomerRepository returns the OUTPUT port implementation.
// Empty DatabaseURL -> in-memory, otherwise PostgreSQL.
// Callers only see port.CustomerRepository, never the concrete type.
func NewCustomerRepository(cfg Config) (port.CustomerRepository, error) {
	if cfg.DatabaseURL == "" {
		return memory.NewCustomerRepository(), nil
	}

	repo, err := postgres.NewCustomerRepository(cfg.DatabaseURL)
	if err != nil {
		return nil, err // explicit, so we never return a nil pointer wrapped in an interface
	}
	return repo, nil
}
