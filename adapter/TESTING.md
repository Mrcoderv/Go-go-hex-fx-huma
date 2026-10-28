# How Mock Tests Work in This Project

This project uses **mock-based unit testing** to verify each layer in isolation, without depending on real databases, HTTP servers, or other external systems.

## Test Files

| File | Package | What it tests | Mock used |
|------|---------|---------------|-----------|
| `internal/core/customer_service_test.go` | `core` | Core business logic | `mockRepository` (output port) |
| `adapter/customer_api_test.go` | `adapter` | HTTP handlers | `mockService` (input port) |
| `adapter/customer_repository_test.go` | `adapter` | In-memory repository | None (uses real implementation) |

## How the Mocks Work

### Core Service Test (`internal/core/customer_service_test.go`)

The core service (`CustomerService`) depends on the **output port** `port.CustomerRepository`. To test it without a real database, we create a fake repository:

```go
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
// ... FindAll, FindByID, Delete follow the same pattern
```

Each test sets fields on the mock to control behavior:

- `TestCreateCustomer` — sets `mockRepository{}`, calls `CreateCustomer`, verifies the repository received the customer and returned it with an ID.
- `TestCreateCustomerRepositoryError` — sets `saveErr: errors.New("db down")`, verifies the error propagates.
- `TestGetCustomerByIDNotFound` — sets `findErr: port.ErrNotFound`, verifies the core returns `ErrNotFound`.
- `TestDeleteCustomerError` — sets `deleteErr: port.ErrNotFound`, verifies the error propagates.

### HTTP Handler Test (`adapter/customer_api_test.go`)

The HTTP handler (`CustomerAPI`) depends on the **input port** `port.CustomerService`. To test it without running an HTTP server, we create a fake service:

```go
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
```

Each test controls the mock's behavior:

- `TestCreateCustomer` — calls `CreateCustomer`, verifies the handler passes the customer to the service and returns the result.
- `TestCreateCustomerServiceError` — sets `createErr`, verifies the error propagates.
- `TestGetCustomerByIDNotFound` — sets `getByIDErr: port.ErrNotFound`, verifies the handler returns an error.
- `TestDeleteCustomerNotFound` — sets `deleteErr: port.ErrNotFound`, verifies the error propagates.

### In-Memory Repository Test (`adapter/customer_repository_test.go`)

This tests the **real** `InMemoryCustomerRepository` (no mock needed — it's the adapter itself):

- `TestInMemorySave` — verifies `Save` assigns ID 1.
- `TestInMemorySaveIncrementsID` — verifies each `Save` increments the ID.
- `TestInMemoryFindAllCopy` — verifies `FindAll` returns a copy (modifying one slice doesn't affect another).
- `TestInMemoryFindByIDNotFound` — verifies `ErrNotFound` for missing IDs.
- `TestInMemoryDelete` / `TestInMemoryDeleteNotFound` — verifies delete behavior.

## Why This Approach Matters

Hexagonal architecture separates layers through **ports (interfaces)**. Mocks let us test each layer against those interfaces:

```
TEST                  MOCK              TESTS
─────────────────────────────────────────────────────────
core test             mockRepository    → CustomerService (core)
handler test          mockService       → CustomerAPI (adapter)
repository test       (real impl)       → InMemoryCustomerRepository (adapter)
```

This means:
- **Core logic is tested without a database.**
- **HTTP handlers are tested without a server.**
- **Each layer can be tested independently.**
- **Swapping the repository (e.g., to Postgres) requires no test changes to core or handlers.**

## Running the Tests

```bash
# Run all tests across all packages
go test ./...

# Run with verbose output to see each test case
go test ./... -v

# Run tests for one package only
go test ./internal/core/... -v
go test ./adapter/... -v
```

Expected output (all passing):

```
=== RUN   TestCreateCustomer
--- PASS: TestCreateCustomer (0.00s)
=== RUN   TestCreateCustomerServiceError
--- PASS: TestCreateCustomerServiceError (0.00s)
...
PASS
ok      customer-management/adapter     0.005s
ok      customer-management/internal/core   0.003s
```

## Adding a New Mock

To test a new component:

1. Identify the port (interface) the component depends on.
2. Create a mock struct that implements that interface.
3. Add fields to control return values and errors.
4. Write test functions that set mock fields, call the component, and assert results.
5. Keep mocks simple — no external dependencies, no file I/O, no network calls.
