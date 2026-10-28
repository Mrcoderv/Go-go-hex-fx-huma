package domain

import "errors"

// Business errors. Every layer can use them, and the HTTP adapter
// turns them into proper status codes (404, 409).
var (
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
