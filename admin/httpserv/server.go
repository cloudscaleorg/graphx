package httpserv

import (
	"net/http"

	"github.com/cloudscaleorg/graphx/admin"
)

// Server implements administration of a GraphX cluster via HTTP REST
// semantics.
//
// Server keeps a reference to an admin.Admin structure
type Server struct {
	*http.Server
	Admin *admin.Admin
}

// New is a constructor for our Server
func New(addr string, a *admin.Admin) *Server {
	mux := http.NewServeMux()
	s := &Server{
		Admin: a,
	}
	s.Register(mux)
	s.Server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return s
}

// Register adds all necessary routes to the provided mux.
func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/charts", ChartCRUD(s.Admin))
	mux.HandleFunc("/api/v1/datasources", DataSourceCRUD(s.Admin))
	mux.HandleFunc("/api/v1/backends", BackendCRUD(s.Admin))
}
