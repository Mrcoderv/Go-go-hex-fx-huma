package port

import (
	"errors"
	"customer-management/internal/domain"
)

// ErrNotFound is returned when a customer does not exist.
var ErrNotFound = errors.New("customer not found")

// CustomerService is the INPUT port.
// It defines the operations the HTTP handlers can ask the core to perform.
// The core service (internal/core) implements this interface.
// Handlers depend only on this interface, never on the concrete service.

type CustomerService interface {
	CreateCustomer(customer domain.Customer) (domain.Customer, error)
	GetCustomers() ([]domain.Customer, error)
	GetCustomerByID(id int) (domain.Customer, error)
	DeleteCustomer(id int) error
}
