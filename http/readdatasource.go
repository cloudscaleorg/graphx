package http

import (
	"encoding/json"
	h "net/http"

	"github.com/cloudscaleorg/graphx/admin"
	fw "github.com/ldelossa/goframework/http"
)

func ReadDataSource(admin admin.DataSource) h.HandlerFunc {
	return func(w h.ResponseWriter, r *h.Request) {
		if r.Method != h.MethodGet {
			resp := fw.NewResponse(fw.CodeMethodNotImplemented, "endpoint only supports GET")
			fw.JError(w, resp, h.StatusNotImplemented)
			return
		}

		ds, err := admin.ReadDataSource()
		if err != nil {
			resp := fw.NewResponse(fw.CodeInternalServerError, "backing store was unavailable")
			fw.JError(w, resp, h.StatusBadRequest)
			return
		}

		err = json.NewEncoder(w).Encode(&ds)
		if err != nil {
			w.WriteHeader(h.StatusInternalServerError)
			return
		}

		return
	}
}
