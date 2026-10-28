package adapter

import (
	"errors"
	"testing"

	"customer-management/internal/domain"
	"customer-management/port"
)

func TestInMemorySave(t *testing.T) {
	repo := NewInMemoryCustomerRepository()

	customer := domain.Customer{Name: "Alice", Email: "alice@example.com", Phone: "123"}
	created, err := repo.Save(customer)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}
	if created.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", created.Name)
	}
}

func TestInMemorySaveIncrementsID(t *testing.T) {
	repo := NewInMemoryCustomerRepository()

	c1, _ := repo.Save(domain.Customer{Name: "Alice"})
	c2, _ := repo.Save(domain.Customer{Name: "Bob"})
	if c1.ID != 1 {
		t.Errorf("first customer should have ID 1, got %d", c1.ID)
	}
	if c2.ID != 2 {
		t.Errorf("second customer should have ID 2, got %d", c2.ID)
	}
}

func TestInMemoryFindAll(t *testing.T) {
	repo := NewInMemoryCustomerRepository()
	repo.Save(domain.Customer{Name: "Alice"})
	repo.Save(domain.Customer{Name: "Bob"})

	customers, err := repo.FindAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(customers) != 2 {
		t.Errorf("expected 2 customers, got %d", len(customers))
	}
}

func TestInMemoryFindAllCopy(t *testing.T) {
	repo := NewInMemoryCustomerRepository()
	repo.Save(domain.Customer{Name: "Alice"})

	first, _ := repo.FindAll()
	second, _ := repo.FindAll()

	// Modifying one slice should not affect the other
	first[0].Name = "Changed"
	if second[0].Name != "Alice" {
		t.Error("FindAll should return a copy, modifications should not leak")
	}
}

func TestInMemoryFindByIDFound(t *testing.T) {
	repo := NewInMemoryCustomerRepository()
	created, _ := repo.Save(domain.Customer{Name: "Alice"})

	found, err := repo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, found.ID)
	}
	if found.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", found.Name)
	}
}

func TestInMemoryFindByIDNotFound(t *testing.T) {
	repo := NewInMemoryCustomerRepository()

	_, err := repo.FindByID(999)
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	if !errors.Is(err, port.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestInMemoryDelete(t *testing.T) {
	repo := NewInMemoryCustomerRepository()
	created, _ := repo.Save(domain.Customer{Name: "Alice"})

	err := repo.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	customers, _ := repo.FindAll()
	if len(customers) != 0 {
		t.Errorf("expected 0 customers after delete, got %d", len(customers))
	}
}

func TestInMemoryDeleteNotFound(t *testing.T) {
	repo := NewInMemoryCustomerRepository()

	err := repo.Delete(999)
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	if !errors.Is(err, port.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
