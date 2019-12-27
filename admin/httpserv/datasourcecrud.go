package httpserv

import (
	"net/http"
	h "net/http"

	"github.com/cloudscaleorg/graphx/admin"
)

func DataSourceCRUD(admin *admin.Admin) h.HandlerFunc {
	return func(w h.ResponseWriter, r *h.Request) {
		switch r.Method {
		case http.MethodGet:
			ReadDataSource(admin).ServeHTTP(w, r)
			return
		case http.MethodPost:
			CreateDataSource(admin).ServeHTTP(w, r)
			return
		case http.MethodPut:
			UpdateDataSource(admin).ServeHTTP(w, r)
			return
		case http.MethodDelete:
			DeleteDataSource(admin).ServeHTTP(w, r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, "method not supported", http.StatusMethodNotAllowed)
			return
		}
	}
}
