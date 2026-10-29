package incoming

import (
	"context"
	"errors"
	"testing"

	"customer-management/internal/domain"
	"customer-management/port"
)

// mockService is a fake implementation of the input port
// used to test the HTTP handlers in isolation.
type mockService struct {
	createCustomer domain.Customer
	createErr      error
	getAll         []domain.Customer
	getAllErr      error
	getByID        domain.Customer
	getByIDErr     error
	deleteID       int
	deleteErr      error
}

func (m *mockService) CreateCustomer(customer domain.Customer) (domain.Customer, error) {
	m.createCustomer = customer
	if m.createErr != nil {
		return domain.Customer{}, m.createErr
	}
	customer.ID = 1
	return customer, nil
}

func (m *mockService) GetCustomers() ([]domain.Customer, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.getAll, nil
}

func (m *mockService) GetCustomerByID(id int) (domain.Customer, error) {
	if m.getByIDErr != nil {
		return domain.Customer{}, m.getByIDErr
	}
	return m.getByID, nil
}

func (m *mockService) UpdateCustomer(id int, customer domain.Customer) (domain.Customer, error) {
	customer.ID = id
	return customer, nil
}

func (m *mockService) DeleteCustomer(id int) error {
	m.deleteID = id
	return m.deleteErr
}

func TestCreateCustomer(t *testing.T) {
	svc := &mockService{}
	api := NewCustomerAPI(svc)

	input := &CreateCustomerInput{}
	input.Body.Name = "Alice"
	input.Body.Email = "alice@example.com"
	input.Body.Phone = "123"

	output, err := api.CreateCustomer(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Body.ID != 1 {
		t.Errorf("expected ID 1, got %d", output.Body.ID)
	}
	if output.Body.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", output.Body.Name)
	}
	if svc.createCustomer.Name != "Alice" {
		t.Errorf("service should receive customer Alice, got %s", svc.createCustomer.Name)
	}
}

func TestCreateCustomerServiceError(t *testing.T) {
	svc := &mockService{createErr: errors.New("service down")}
	api := NewCustomerAPI(svc)

	_, err := api.CreateCustomer(context.Background(), &CreateCustomerInput{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListCustomers(t *testing.T) {
	expected := []domain.Customer{{ID: 1, Name: "Alice"}}
	svc := &mockService{getAll: expected}
	api := NewCustomerAPI(svc)

	output, err := api.ListCustomers(context.Background(), &struct{}{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Body.Customers) != 1 {
		t.Errorf("expected 1 customer, got %d", len(output.Body.Customers))
	}
	if output.Body.Customers[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", output.Body.Customers[0].Name)
	}
}

func TestListCustomersError(t *testing.T) {
	svc := &mockService{getAllErr: errors.New("db down")}
	api := NewCustomerAPI(svc)

	_, err := api.ListCustomers(context.Background(), &struct{}{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetCustomerByIDFound(t *testing.T) {  // random data to test the handler return false
	svc := &mockService{getByID: domain.Customer{ID: 7, Name: "Bob"}}
	api := NewCustomerAPI(svc)

	output, err := api.GetCustomerByID(context.Background(), &GetCustomerInput{ID: 7})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Body.ID != 7 {
		t.Errorf("expected ID 7, got %d", output.Body.ID)
	}
	if svc.getByIDErr != nil {
		t.Errorf("unexpected error set on mock")
	}
}

func TestGetCustomerByIDNotFound(t *testing.T) {
	svc := &mockService{getByIDErr: port.ErrNotFound}
	api := NewCustomerAPI(svc)

	_, err := api.GetCustomerByID(context.Background(), &GetCustomerInput{ID: 999})
	if err == nil {
		t.Fatal("expected 404 error, got nil")
	}
	if err.Error() != "customer not found" {
		t.Errorf("expected 'customer not found', got %v", err)
	}
}

func TestDeleteCustomer(t *testing.T) {  // random deletingt the dta of the id 3 .
	svc := &mockService{}
	api := NewCustomerAPI(svc)

	output, err := api.DeleteCustomer(context.Background(), &DeleteCustomerInput{ID: 3})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Body.Message != "customer deleted" {
		t.Errorf("expected 'customer deleted', got %s", output.Body.Message)
	}
	if svc.deleteID != 3 {
		t.Errorf("expected deleteID 3, got %d", svc.deleteID)
	}
}

func TestDeleteCustomerNotFound(t *testing.T) {
	svc := &mockService{deleteErr: port.ErrNotFound}
	api := NewCustomerAPI(svc)

	_, err := api.DeleteCustomer(context.Background(), &DeleteCustomerInput{ID: 999})
	if err == nil {
		t.Fatal("expected 404 error, got nil")
	}
	if err.Error() != "customer not found" {
		t.Errorf("expected 'customer not found', got %v", err)
	}
}
