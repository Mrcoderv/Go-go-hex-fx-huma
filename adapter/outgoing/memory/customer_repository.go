package memory

import (
	"customer-management/internal/domain"
	"customer-management/port"
)

// CustomerRepository stores customers in a slice.
// It implements the output port (port.CustomerRepository).
type CustomerRepository struct {
	customers []domain.Customer
	nextID    int
}

// NewCustomerRepository creates an empty in-memory repository.
func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{
		customers: make([]domain.Customer, 0),
		nextID:    1,
	}
}

// Save assigns a new ID and appends the customer.
func (r *CustomerRepository) Save(
	customer domain.Customer,
) (domain.Customer, error) {
	customer.ID = r.nextID
	r.nextID++
	r.customers = append(r.customers, customer)
	return customer, nil
}

// FindAll returns a copy of all customers.
func (r *CustomerRepository) FindAll() ([]domain.Customer, error) {
	customers := make([]domain.Customer, len(r.customers))
	copy(customers, r.customers)
	return customers, nil
}

// FindByID returns the customer with the given ID, or ErrNotFound.
func (r *CustomerRepository) FindByID(id int) (domain.Customer, error) {
	for _, c := range r.customers {
		if c.ID == id {
			return c, nil
		}
	}
	return domain.Customer{}, port.ErrNotFound
}

// Delete removes the customer with the given ID, or returns ErrNotFound.
func (r *CustomerRepository) Delete(id int) error {
	for i, c := range r.customers {
		if c.ID == id {
			r.customers = append(r.customers[:i], r.customers[i+1:]...)
			return nil
		}
	}
	return port.ErrNotFound
}

// Compile-time check: this adapter must satisfy the output port.
var _ port.CustomerRepository = (*CustomerRepository)(nil)
