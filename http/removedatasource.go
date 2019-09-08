package http

import (
	"encoding/json"
	h "net/http"

	"github.com/cloudscaleorg/graphx"
	"github.com/cloudscaleorg/graphx/admin"
	fw "github.com/ldelossa/goframework/http"
)

func RemoveDataSource(a admin.DataSource) h.HandlerFunc {
	return func(w h.ResponseWriter, r *h.Request) {
		if r.Method != h.MethodPut {
			resp := fw.NewResponse(fw.CodeMethodNotImplemented, "endpoint only supports PUT")
			fw.JError(w, resp, h.StatusNotImplemented)
			return
		}

		var ds graphx.DataSource
		err := json.NewDecoder(r.Body).Decode(&ds)
		if err != nil {
			resp := fw.NewResponse(fw.CodeFailedSerialization, "could not validate provided json")
			fw.JError(w, resp, h.StatusBadRequest)
			return
		}

		err = a.DeleteDataSource(&ds)
		if err != nil {
			switch {
			case (err == admin.ErrNotFound):
				resp := fw.NewResponse(fw.CodeNotFound, "resource being updated not found")
				fw.JError(w, resp, h.StatusNotFound)
				return
			case (err == admin.ErrStore{}):
				resp := fw.NewResponse(fw.CodeInternalServerError, "an internal error occured")
				fw.JError(w, resp, h.StatusNotFound)
				return
			}
		}

		w.WriteHeader(h.StatusOK)
		return
	}
}
