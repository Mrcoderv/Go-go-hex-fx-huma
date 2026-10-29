package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"customer-management/adapter/incoming"
	"customer-management/adapter/outgoing"
	"customer-management/internal/core"
	"customer-management/port"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

func NewMux() *http.ServeMux {
	return http.NewServeMux()
}

// NewCustomerRepository gives Fx the OUTPUT port. main.go does not know
// which storage is behind it; the outgoing layer decides from DATABASE_URL.
func NewCustomerRepository() (port.CustomerRepository, error) {
	return outgoing.NewCustomerRepository(outgoing.Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	})
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
	handler *incoming.CustomerAPI,
) {
	incoming.RegisterCustomerRoutes(api, handler)
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
			NewCustomerRepository,
			core.NewCustomerService,
			ProvideCustomerService,
			incoming.NewCustomerAPI,
			incoming.NewHumaAPI,
		),
		fx.Invoke(
			RegisterRoutes,
			StartHTTPServer,
		),
	).Run()
}
