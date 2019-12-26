package http

import (
	"encoding/json"
	h "net/http"

	"github.com/cloudscaleorg/graphx/admin"
	fw "github.com/ldelossa/goframework/http"
	"github.com/rs/zerolog/log"
)

func ReadDataSource(admin *admin.Admin) h.HandlerFunc {
	logger := log.With().Str("component", "ReadDataSourceHandler").Logger()
	return func(w h.ResponseWriter, r *h.Request) {
		if r.Method != h.MethodGet {
			resp := fw.NewResponse(fw.CodeMethodNotImplemented, "endpoint only supports GET")
			fw.JError(w, resp, h.StatusNotImplemented)
			return
		}
		v, err := admin.ReadDataSource()
		if err != nil {
			logger.Error().Msgf("failed to read datasources from admin interface: %v", err)
			resp := fw.NewResponse(fw.CodeInternalServerError, "backing store was unavailable")
			fw.JError(w, resp, h.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(&v)
		if err != nil {
			logger.Error().Msgf("failed to deserialize datasource: %v", err)
			w.WriteHeader(h.StatusInternalServerError)
			return
		}
		return
	}
}
