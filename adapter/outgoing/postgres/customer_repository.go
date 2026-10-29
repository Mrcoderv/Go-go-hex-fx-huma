package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"customer-management/internal/domain"
	"customer-management/port"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const schema = `CREATE TABLE IF NOT EXISTS customers (
	id    SERIAL PRIMARY KEY,
	name  TEXT NOT NULL,
	email TEXT NOT NULL,
	phone TEXT NOT NULL DEFAULT ''
)`

// queryTimeout bounds every database call. The output port has no
// context parameter yet, so the adapter creates its own.
const queryTimeout = 5 * time.Second

// CustomerRepository saves customers in PostgreSQL.
// It implements the output port (port.CustomerRepository).
type CustomerRepository struct {
	pool *pgxpool.Pool
}

// NewCustomerRepository connects to PostgreSQL and creates the table if missing.
func NewCustomerRepository(dsn string) (*CustomerRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if _, err := pool.Exec(ctx, schema); err != nil {
		pool.Close()
		return nil, err
	}
	return &CustomerRepository{pool: pool}, nil
}

// Save inserts the customer and returns it with the generated ID.
func (r *CustomerRepository) Save(customer domain.Customer) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := r.pool.QueryRow(ctx,
		`INSERT INTO customers (name, email, phone)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		customer.Name, customer.Email, customer.Phone,
	).Scan(&customer.ID)
	fmt.Println("customer id is created ", customer.ID)
	if err != nil {
		return domain.Customer{}, err
	}
	return customer, nil
}

// FindAll returns every customer ordered by ID. Never returns a nil slice,
// so the JSON response is [] and not null (same as the in-memory adapter).
func (r *CustomerRepository) FindAll() ([]domain.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, email, phone FROM customers ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]domain.Customer, 0)
	for rows.Next() {
		var c domain.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, rows.Err()
}

// FindByID returns the customer with the given ID, or port.ErrNotFound.
func (r *CustomerRepository) FindByID(id int) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var c domain.Customer
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, email, phone FROM customers WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Email, &c.Phone)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Customer{}, port.ErrNotFound
	}
	if err != nil {
		return domain.Customer{}, err
	}
	return c, nil
}

// Delete removes the customer with the given ID, or returns port.ErrNotFound.
func (r *CustomerRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	tag, err := r.pool.Exec(ctx, `DELETE FROM customers WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return port.ErrNotFound
	}
	return nil
}

// Compile-time check: this adapter must satisfy the output port.
var _ port.CustomerRepository = (*CustomerRepository)(nil)
