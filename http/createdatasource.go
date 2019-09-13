package http

import (
	"encoding/json"
	h "net/http"

	"github.com/cloudscaleorg/graphx"
	"github.com/cloudscaleorg/graphx/admin"
	fw "github.com/ldelossa/goframework/http"
)

func CreateDataSource(admin admin.DataSource) h.HandlerFunc {
	return func(w h.ResponseWriter, r *h.Request) {
		if r.Method != h.MethodPost {
			resp := fw.NewResponse(fw.CodeMethodNotImplemented, "endpoint only supports POST")
			fw.JError(w, resp, h.StatusNotImplemented)
			return
		}

		var v graphx.DataSource
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			resp := fw.NewResponse(fw.CodeFailedSerialization, "could not validate provided json")
			fw.JError(w, resp, h.StatusBadRequest)
			return
		}

		resp := fw.NewResponse(fw.CodeSuccess, "DataSource stored")
		err := json.NewEncoder(w).Encode(&resp)
		if err != nil {
			w.WriteHeader(h.StatusInternalServerError)
			return
		}

		return
	}
}
