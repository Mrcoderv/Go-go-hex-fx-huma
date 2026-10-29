package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"customer-management/internal/domain"
	"customer-management/port"
)

// These are integration tests. They need a real database and are skipped
// unless TEST_DATABASE_URL is set, e.g.
//
//	TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/customers?sslmode=disable go test ./adapter/outgoing/postgres/...
func newTestRepo(t *testing.T) *CustomerRepository {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping postgres integration test")
	}

	repo, err := NewCustomerRepository(dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(repo.pool.Close)

	if _, err := repo.pool.Exec(context.Background(), `TRUNCATE customers RESTART IDENTITY`); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	return repo
}

func TestPostgresSaveAndFindByID(t *testing.T) {
	repo := newTestRepo(t)

	saved, err := repo.Save(domain.Customer{Name: "Alice", Email: "alice@example.com", Phone: "123"})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if saved.ID != 1 {
		t.Errorf("expected ID 1, got %d", saved.ID)
	}

	found, err := repo.FindByID(saved.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if found != saved {
		t.Errorf("expected %+v, got %+v", saved, found)
	}
}

func TestPostgresFindAll(t *testing.T) {
	repo := newTestRepo(t)

	all, err := repo.FindAll()
	if err != nil || all == nil || len(all) != 0 {
		t.Fatalf("expected empty non-nil slice, got %v, err %v", all, err)
	}

	repo.Save(domain.Customer{Name: "Alice"})
	repo.Save(domain.Customer{Name: "Bob"})

	all, err = repo.FindAll()
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if len(all) != 2 || all[0].Name != "Alice" || all[1].Name != "Bob" {
		t.Errorf("unexpected result: %+v", all)
	}
}

func TestPostgresFindByIDNotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.FindByID(999)
	if !errors.Is(err, port.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresDelete(t *testing.T) {
	repo := newTestRepo(t)

	saved, _ := repo.Save(domain.Customer{Name: "Alice"})
	if err := repo.Delete(saved.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.FindByID(saved.ID); !errors.Is(err, port.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestPostgresDeleteNotFound(t *testing.T) {
	repo := newTestRepo(t)

	if err := repo.Delete(999); !errors.Is(err, port.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
