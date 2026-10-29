package incoming

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func NewHumaAPI(mux *http.ServeMux) huma.API {
	config := huma.DefaultConfig("Customer Management API", "1.0.0")
	return humago.New(mux, config)
}
