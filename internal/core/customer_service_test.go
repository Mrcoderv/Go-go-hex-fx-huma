package core

import (
	"errors"
	"testing"

	"customer-management/internal/domain"
	"customer-management/port"
)

// mockRepository is a fake implementation of the output port
// used to test the core service in isolation.
type mockRepository struct {
	saved     domain.Customer
	all       []domain.Customer
	findByID  domain.Customer
	deleteID  int
	findErr   error
	deleteErr error
	saveErr   error
	findAllErr error
}

func (m *mockRepository) Save(customer domain.Customer) (domain.Customer, error) {
	m.saved = customer
	if m.saveErr != nil {
		return domain.Customer{}, m.saveErr
	}
	customer.ID = 1
	return customer, nil
}

func (m *mockRepository) FindAll() ([]domain.Customer, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	return m.all, nil
}

func (m *mockRepository) FindByID(id int) (domain.Customer, error) {
	if m.findErr != nil {
		return domain.Customer{}, m.findErr
	}
	return m.findByID, nil
}

func (m *mockRepository) Update(id int, customer domain.Customer) (domain.Customer, error) {
	customer.ID = id
	return customer, nil
}

func (m *mockRepository) Delete(id int) error {
	m.deleteID = id
	return m.deleteErr
}

func TestCreateCustomer(t *testing.T) {
	repo := &mockRepository{}
	service := NewCustomerService(repo)

	customer := domain.Customer{Name: "Alice", Email: "alice@example.com", Phone: "123"}
	created, err := service.CreateCustomer(customer)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}
	if created.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", created.Name)
	}
	if repo.saved.Name != "Alice" {
		t.Errorf("repository should receive the customer, got %s", repo.saved.Name)
	}
}

func TestCreateCustomerRepositoryError(t *testing.T) {
	repo := &mockRepository{saveErr: errors.New("db down")}
	service := NewCustomerService(repo)

	_, err := service.CreateCustomer(domain.Customer{Name: "Bob"})
	if err == nil {
		t.Fatal("expected error from repository, got nil")
	}
	if err.Error() != "db down" {
		t.Errorf("expected 'db down', got %v", err)
	}
}

func TestGetCustomers(t *testing.T) {
	expected := []domain.Customer{{ID: 1, Name: "Alice"}}
	repo := &mockRepository{all: expected}
	service := NewCustomerService(repo)

	customers, err := service.GetCustomers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(customers) != 1 {
		t.Errorf("expected 1 customer, got %d", len(customers))
	}
	if customers[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", customers[0].Name)
	}
}

func TestGetCustomersError(t *testing.T) {
	repo := &mockRepository{findAllErr: errors.New("db down")}
	service := NewCustomerService(repo)

	_, err := service.GetCustomers()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetCustomerByID(t *testing.T) {
	expected := domain.Customer{ID: 5, Name: "Alice"}
	repo := &mockRepository{findByID: expected}
	service := NewCustomerService(repo)

	customer, err := service.GetCustomerByID(5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if customer.ID != 5 {
		t.Errorf("expected ID 5, got %d", customer.ID)
	}
}

func TestGetCustomerByIDNotFound(t *testing.T) {
	repo := &mockRepository{findErr: port.ErrNotFound}
	service := NewCustomerService(repo)

	_, err := service.GetCustomerByID(999)
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	if !errors.Is(err, port.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteCustomer(t *testing.T) {
	repo := &mockRepository{}
	service := NewCustomerService(repo)

	err := service.DeleteCustomer(3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.deleteID != 3 {
		t.Errorf("expected deleteID 3, got %d", repo.deleteID)
	}
}

func TestDeleteCustomerError(t *testing.T) {
	repo := &mockRepository{deleteErr: port.ErrNotFound}
	service := NewCustomerService(repo)

	err := service.DeleteCustomer(999)
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	if !errors.Is(err, port.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
