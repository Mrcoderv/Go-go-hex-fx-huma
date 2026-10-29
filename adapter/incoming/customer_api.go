package incoming

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"customer-management/internal/domain"
	"customer-management/port"

	"github.com/danielgtaylor/huma/v2"
)

// CustomerAPI handles HTTP requests for customers.
// It depends on the INPUT port (port.CustomerService), not the concrete core.
type CustomerAPI struct {
	service port.CustomerService
}

// NewCustomerAPI creates a new customer API handler.
func NewCustomerAPI(service port.CustomerService) *CustomerAPI {
	return &CustomerAPI{service: service}
}

// -------- POST /customers --------

type CreateCustomerInput struct {
	Body struct {
		Name  string `json:"name" minLength:"1" doc:"Customer name"`
		Email string `json:"email" format:"email" doc:"Customer email"`
		Phone string `json:"phone" doc:"Customer phone"`
	}
}

type CreateCustomerOutput struct {
	Body domain.Customer
}

func (h *CustomerAPI) CreateCustomer(
	ctx context.Context,
	input *CreateCustomerInput,
) (*CreateCustomerOutput, error) {
	customer := domain.Customer{
		Name:  input.Body.Name,
		Email: input.Body.Email,
		Phone: input.Body.Phone,
	}

	created, err := h.service.CreateCustomer(customer)
	if err != nil {
		return nil, err
	}

	return &CreateCustomerOutput{Body: created}, nil
}

// -------- GET /customers --------

type ListCustomersOutput struct {
	Body struct {
		Customers []domain.Customer `json:"customers"`
	}
}

func (h *CustomerAPI) ListCustomers(
	ctx context.Context,
	input *struct{},
) (*ListCustomersOutput, error) {
	customers, err := h.service.GetCustomers()
	if err != nil {
		return nil, err
	}
	return &ListCustomersOutput{
		Body: struct {
			Customers []domain.Customer `json:"customers"`
		}{Customers: customers},
	}, nil
}

// -------- GET /customers/{id} --------

type GetCustomerInput struct {
	ID int `path:"id" doc:"Customer ID"`
}

type GetCustomerOutput struct {
	Body domain.Customer
}

func (h *CustomerAPI) GetCustomerByID(
	ctx context.Context,
	input *GetCustomerInput,
) (*GetCustomerOutput, error) {
	customer, err := h.service.GetCustomerByID(input.ID)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "customer not found")
		}
		return nil, err
	}
	return &GetCustomerOutput{Body: customer}, nil
}

// -------- DELETE /customers/{id} --------

type DeleteCustomerInput struct {
	ID int `path:"id" doc:"Customer ID"`
}

type DeleteCustomerOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func (h *CustomerAPI) DeleteCustomer(
	ctx context.Context,
	input *DeleteCustomerInput,
) (*DeleteCustomerOutput, error) {
	fmt.Println("Deleting customer with ID:", input.ID)
	err := h.service.DeleteCustomer(input.ID)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "customer not found")
		}
		return nil, err
	}
	return &DeleteCustomerOutput{
		Body: struct {
			Message string `json:"message"`
		}{Message: "customer deleted"},
	}, nil
}

// -------- Route registration --------

func RegisterCustomerRoutes(api huma.API, handler *CustomerAPI) {
	huma.Register(api, huma.Operation{
		OperationID: "create-customer",
		Method:      "POST",
		Path:        "/customers",
		Summary:     "Create a customer",
	}, handler.CreateCustomer)

	huma.Register(api, huma.Operation{
		OperationID: "list-customers",
		Method:      "GET",
		Path:        "/customers",
		Summary:     "List all customers",
	}, handler.ListCustomers)

	huma.Register(api, huma.Operation{
		OperationID: "get-customer",
		Method:      "GET",
		Path:        "/customers/{id}",
		Summary:     "Get a customer by ID",
	}, handler.GetCustomerByID)

	huma.Register(api, huma.Operation{
		OperationID: "delete-customer",
		Method:      "DELETE",
		Path:        "/customers/{id}",
		Summary:     "Delete a customer",
	}, handler.DeleteCustomer)
	huma.Register(api, huma.Operation{
		OperationID: "put-customer",
		Method:      "PUT",
		Path:        "/customers/{id}",
		Summary:     "Update a customer",
	}, handler.UpdateCustomer)

}

// -------- PUT /customers/{id} --------

type UpdateCustomerInput struct {
	ID   int `path:"id" doc:"Customer ID"`
	Body struct {
		Name  string `json:"name" minLength:"1" doc:"Customer name"`
		Email string `json:"email" format:"email" doc:"Customer email"`
		Phone string `json:"phone" doc:"Customer phone"`
	}
}

type UpdateCustomerOutput struct {
	Body domain.Customer
}

func (h *CustomerAPI) UpdateCustomer(
	ctx context.Context,
	input *UpdateCustomerInput,
) (*UpdateCustomerOutput, error) {
	customer, err := h.service.UpdateCustomer(input.ID, domain.Customer{
		Name:  input.Body.Name,
		Email: input.Body.Email,
		Phone: input.Body.Phone,
	})
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "customer not found")
		}
		return nil, err
	}
	return &UpdateCustomerOutput{Body: customer}, nil
}
