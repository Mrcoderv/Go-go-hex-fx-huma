package core

import (
	"customer-management/internal/domain"
	"customer-management/port"
)

// CustomerService is the core business logic.
// It implements the input port (port.CustomerService) and
// uses the output port (port.CustomerRepository) to persist data.
type CustomerService struct {
	repository port.CustomerRepository
}

// NewCustomerService creates a new core service wired to a repository.
func NewCustomerService(repository port.CustomerRepository) *CustomerService {
	return &CustomerService{repository: repository}
}

// CreateCustomer delegates saving to the repository. 
// using the input port interface, it does not depend on the concrete repository implementation.
func (s *CustomerService) CreateCustomer(
	customer domain.Customer,
) (domain.Customer, error) {
	return s.repository.Save(customer)
}

// GetCustomers returns every customer. 
//
func (s *CustomerService) GetCustomers() ([]domain.Customer, error) {
	return s.repository.FindAll()
}

// GetCustomerByID returns a single customer by its ID.
func (s *CustomerService) GetCustomerByID(id int) (domain.Customer, error) {
	return s.repository.FindByID(id)
}

// DeleteCustomer removes a customer by its ID.
func (s *CustomerService) DeleteCustomer(id int) error {
	return s.repository.Delete(id)
}
