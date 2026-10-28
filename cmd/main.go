package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"customer-management/adapter"
	"customer-management/internal/core"
	"customer-management/port"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

func NewMux() *http.ServeMux {
	return http.NewServeMux()
}

// ProvideCustomerRepository tells Fx how to turn the concrete
// in-memory repository into the output port interface.
func ProvideCustomerRepository(
	repo *adapter.InMemoryCustomerRepository,
) port.CustomerRepository {
	return repo
}

// ProvideCustomerService tells Fx how to turn the concrete
// core service into the input port interface.
func ProvideCustomerService(
	service *core.CustomerService,
) port.CustomerService {
	return service
}

func RegisterRoutes(
	api huma.API,
	handler *adapter.CustomerAPI,
) {
	adapter.RegisterCustomerRoutes(api, handler)
}

func StartHTTPServer(
	lc fx.Lifecycle,
	mux *http.ServeMux,
) {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}

			go func() {
				if err := server.Serve(listener); err != nil &&
					err != http.ErrServerClosed {
					log.Printf("HTTP server error: %v", err)
				}
			}()

			log.Println("Server running at http://localhost:8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func main() {
	fx.New(
		fx.Provide(
			NewMux,
			adapter.NewInMemoryCustomerRepository,
			ProvideCustomerRepository,
			core.NewCustomerService,
			ProvideCustomerService,
			adapter.NewCustomerAPI,
			adapter.NewHumaAPI,
		),
		fx.Invoke(
			RegisterRoutes,
			StartHTTPServer,
		),
	).Run()
}
