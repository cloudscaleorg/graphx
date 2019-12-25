package http

import (
	"net/http"
	h "net/http"

	"github.com/cloudscaleorg/graphx/admin"
)

func BackendCRUD(admin admin.Backend) h.HandlerFunc {
	return func(w h.ResponseWriter, r *h.Request) {
		switch r.Method {
		case http.MethodGet:
			ReadBackend(admin).ServeHTTP(w, r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
	}
}
