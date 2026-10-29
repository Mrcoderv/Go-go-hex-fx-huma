package outgoing

import (
	"testing"

	"customer-management/adapter/outgoing/memory"
)

func TestEmptyDatabaseURLUsesMemory(t *testing.T) {
	repo, err := NewCustomerRepository(Config{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := repo.(*memory.CustomerRepository); !ok {
		t.Errorf("expected *memory.CustomerRepository, got %T", repo)
	}
}
