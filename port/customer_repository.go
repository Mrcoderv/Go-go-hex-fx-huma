package port

import "customer-management/internal/domain"

// CustomerRepository is the OUTPUT port.
// It defines how the core saves and reads customers.
// An adapter (in-memory or postgres) implements this interface.
// The core depends only on this interface, never on a concrete database.

// if we implement a new database, we just need to implement this interface and the core will work with it without any changes.
type CustomerRepository interface {
	Save(customer domain.Customer) (domain.Customer, error)
	FindAll() ([]domain.Customer, error)
	FindByID(id int) (domain.Customer, error)
	Delete(id int) error
}
