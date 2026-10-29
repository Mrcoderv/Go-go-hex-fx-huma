# Customer Management API

A simple Customer Management API built with **Go** using **Hexagonal Architecture** (Ports & Adapters) and **Uber Fx** for dependency injection.

## What This Project Does

A REST API for managing customers with four endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/customers` | Create a new customer |
| `GET` | `/customers` | List all customers |
| `GET` | `/customers/{id}` | Get a customer by ID |
| `DELETE` | `/customers/{id}` | Delete a customer by ID |

## Project Structure

```
cms/
├── cmd/                          # Application entry point (composition root)
│   └── main.go                   # Fx app wiring + HTTP server lifecycle
├── internal/
│   ├── domain/                   # Domain entities (plain Go structs)
│   │   └── customer.go           # Customer model: ID, Name, Email, Phone
│   └── core/                     # Business logic (use cases)
│       └── customer_service.go   # CustomerService implements input port
├── port/                         # Ports (interfaces) — the hexagon boundary
│   ├── customer_service.go       # INPUT port: interface handlers call
│   └── customer_repository.go    # OUTPUT port: interface core calls
├── adapter/
│   ├── incoming/                 # INCOMING adapters: handle requests
│   │   ├── customer_api.go       # HTTP handlers (depend on input port)
│   │   └── huma_api.go           # Huma API setup on net/http mux
│   └── outgoing/                 # OUTGOING adapters: save data
│       ├── repository.go         # Factory: picks memory or postgres (main.go only sees the port)
│       ├── memory/               # In-memory repository (implements output port)
│       └── postgres/             # PostgreSQL repository (implements output port)
├── docker-compose.yml            # only for the  PostgreSQL  so that if the native local device db didn't work   .
├── go.mod                        # Module definition + dependencies
├── go.sum                        # Dependency checksums
└── README.md                     # This file
```

## Hexagonal Architecture — How the Ports Connect Everything

Hexagonal architecture isolates the **core business logic** from external concerns (HTTP, database). Communication happens through **ports** (interfaces):

### Input Port (`port.CustomerService`)

Defines what the HTTP handlers can ask the core to do.

```go
// port/customer_service.go
type CustomerService interface {
    CreateCustomer(customer domain.Customer) (domain.Customer, error)
    GetCustomers() ([]domain.Customer, error)
    GetCustomerByID(id int) (domain.Customer, error)
    DeleteCustomer(id int) error
}
```

- **Who implements it**: `internal/core.CustomerService`
- **Who depends on it**: `incoming.CustomerAPI` (HTTP handlers)
- **Direction**: Adapter → Core

### Output Port (`port.CustomerRepository`)

Defines how the core saves and reads customers.

```go
// port/customer_repository.go
type CustomerRepository interface {
    Save(customer domain.Customer) (domain.Customer, error)
    FindAll() ([]domain.Customer, error)
    FindByID(id int) (domain.Customer, error)
    Delete(id int) error
}
```

- **Who implements it**: `memory.CustomerRepository` and `postgres.CustomerRepository` (`adapter/outgoing/`)
- **Who depends on it**: `internal/core.CustomerService`
- **Direction**: Core → Adapter

### Connection Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP CLIENT                          │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP request
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  ADAPTER LAYER (framework-specific)                         │
│  adapter/customer_api.go  ◄── depends on ─── INPUT PORT     │
│  (CustomerAPI handlers)                                       │
└──────────────────────┬──────────────────────────────────────┘
                       │ calls interface methods
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  CORE LAYER (business logic, no framework dependencies)      │
│  internal/core/customer_service.go  ◄── implements ───┐     │
│  (CustomerService)                                        │     │
└──────────────────────┬──────────────────────────────────────┘
                       │ calls interface methods
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  ADAPTER LAYER (persistence)                                 │
│  adapter/customer_repository.go  ◄── implements ─── OUTPUT   │
│  (InMemoryCustomerRepository)         PORT                   │
└──────────────────────┬──────────────────────────────────────┘
                       │ reads/writes
                       ▼
              postgarase

```

## Component Explanations (Simple Language)

### `internal/domain/customer.go` — Domain Model

Plain Go struct representing a customer. Contains only data (ID, Name, Email, Phone). No logic, no dependencies on frameworks.

### `port/customer_service.go` — Input Port

An interface listing the operations the HTTP layer can request from the core. The core service implements this. Handlers depend on **this interface**, not on a concrete struct, so you can swap the core without changing handlers.

### `port/customer_repository.go` — Output Port

An interface listing the operations the core can request from the persistence layer. An adapter (in-memory, postgres, etc.) implements this. The core depends on **this interface**, not on a concrete database, so you can swap storage without changing the core.

### `internal/core/customer_service.go` — Core Service

Contains business logic. It receives a repository (output port) and implements the input port. For now it mostly delegates to the repository, but this is where validation or business rules would go.

### `adapter/customer_repository.go` — In-Memory Repository

Stores customers in a Go slice. Implements the output port. Data is lost when the process stops — good for demos and tests.

### `adapter/customer_api.go` — HTTP Handlers

Handles HTTP requests using Huma. Depends on the input port interface. Maps HTTP requests to core service calls and returns responses. Contains route registration for all four endpoints.

### `adapter/huma_api.go` — Huma Setup

Creates a Huma API instance wrapping Go's standard `http.ServeMux`. Huma provides OpenAPI documentation and request/response validation.

### `cmd/main.go` — Entry Point (Composition Root)

Uses **Uber Fx** to wire everything together:
1. Creates dependencies (mux, repository, service, API, handlers)
2. Resolves interfaces to concrete implementations via provider functions
3. Starts the HTTP server on `:8080`
4. Registers all routes

## How the Data Flows

**Create customer** (`POST /customers`):
1. HTTP request hits `CustomerAPI.CreateCustomer`
2. Handler calls `port.CustomerService.CreateCustomer` (input port)
3. Core service calls `port.CustomerRepository.Save` (output port)
4. In-memory repository assigns an ID and stores the customer
5. Result flows back: repository → core → handler → HTTP response

**Get customer by ID** (`GET /customers/{id}`):
1. Handler extracts `id` from the URL path
2. Handler calls `GetCustomerByID(id)` on the input port
3. Core calls `FindByID(id)` on the output port
4. Repository searches the slice; returns `ErrNotFound` if missing
5. Handler converts `ErrNotFound` to HTTP 404

## Dependencies

| Package | Purpose |
|---------|---------|
| `go.uber.org/fx` | Dependency injection framework |
| `github.com/danielgtaylor/huma/v2` | Web framework + OpenAPI |
| `github.com/danielgtaylor/huma/v2/adapters/humago` | Huma adapter for net/http |
| `github.com/jackc/pgx/v5` | PostgreSQL driver and connection pool |

## Getting Started

### Prerequisites

Go 1.27.0 or newer.

### Install Dependencies

```bash
go mod tidy   # downloads pgx and updates go.sum
```

### Run the Server

```bash
go run ./cmd
```

Server starts at `http://localhost:8080`.

### Test the Endpoints

```bash
# Create a customer
curl -X POST http://localhost:8080/customers \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com","phone":"1234567890"}'

# List all customers
curl http://localhost:8080/customers

# Get a customer by ID
curl http://localhost:8080/customers/1

# Delete a customer by ID
curl -X DELETE http://localhost:8080/customers/1
```

### Build

```bash
go build ./...
```

### Run Tests

```bash
# Run all tests across all packages
go test ./...

# Run with verbose output to see each test case
go test ./... -v
```

Postgres integration tests are skipped unless `TEST_DATABASE_URL` is set:

```bash
docker compose up -d db
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/customers?sslmode=disable" go test ./adapter/outgoing/postgres/... -v
```

For details on how the mock tests work and what each test file covers, see **[TESTING.md](./TESTING.md)**.

## Choosing the Storage (In-Memory or PostgreSQL)

Both repositories implement `port.CustomerRepository`. The factory in `adapter/outgoing/repository.go` picks one at startup using the `DATABASE_URL` environment variable, so `cmd/main.go` never imports `memory` or `postgres`:

| `DATABASE_URL` | Repository used |
|----------------|-----------------|
| not set | In-memory (data is lost on restart) |
| set | PostgreSQL (the `customers` table is created automatically) |

To add another database, create a new folder under `adapter/outgoing/` and add one branch in the factory.

```bash
# In-memory (default)
go run ./cmd

# PostgreSQL
docker compose up -d db
DATABASE_URL="postgres://postgres:postgres@localhost:5432/customers?sslmode=disable" go run ./cmd
```

No changes were needed in `internal/core`, `port`, or the HTTP handlers. The in-memory adapter is kept (it is not replaced) so demos and tests still run without a database.

```bash

mind map 
HTTP Client
    │ HTTP request
    ▼
┌──────────────────────────┐
│  ADAPTER (handlers)      │  ──depends on──►  INPUT PORT (interface)
│  adapter/incoming/        │                    port/customer_service.go
└──────────────────────────┘
    │ calls interface methods
    ▼
┌──────────────────────────┐
│  CORE (service)          │  ──implements──►  INPUT PORT
│  internal/core/          │  ──depends on──►  OUTPUT PORT (interface)
│  customer_service.go     │                    port/customer_repository.go
└──────────────────────────┘
    │ calls interface methods
    ▼
┌──────────────────────────┐
│  ADAPTER (repository)    │  ──implements──►  OUTPUT PORT
│  adapter/outgoing/        │
│  memory | postgres       │
└──────────────────────────┘
    │ reads/writes
    ▼
   In-memory slice  or  PostgreSQL

```
