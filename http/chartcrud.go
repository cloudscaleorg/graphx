package http

import (
	"net/http"
	h "net/http"

	"github.com/cloudscaleorg/graphx/admin"
)

func ChartCRUD(admin admin.Chart) h.HandlerFunc {
	return func(w h.ResponseWriter, r *h.Request) {
		switch r.Method {
		case http.MethodGet:
			ReadChart(admin).ServeHTTP(w, r)
			return
		case http.MethodPost:
			CreateChart(admin).ServeHTTP(w, r)
			return
		case http.MethodPut:
			UpdateChart(admin).ServeHTTP(w, r)
			return
		case http.MethodDelete:
			DeleteChart(admin).ServeHTTP(w, r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
	}
}
