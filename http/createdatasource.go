package http

import (
	"encoding/json"
	h "net/http"

	"github.com/cloudscaleorg/graphx"
	"github.com/cloudscaleorg/graphx/admin"
	fw "github.com/ldelossa/goframework/http"
	"github.com/rs/zerolog/log"
)

func CreateDataSource(admin admin.DataSource) h.HandlerFunc {
	logger := log.With().Str("component", "CreateDataSourceHandler").Logger()
	return func(w h.ResponseWriter, r *h.Request) {
		if r.Method != h.MethodPost {
			resp := fw.NewResponse(fw.CodeMethodNotImplemented, "endpoint only supports POST")
			fw.JError(w, resp, h.StatusNotImplemented)
			return
		}

		var v graphx.DataSource
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			logger.Error().Msgf("failed to deserialize datasource: %v", err)
			resp := fw.NewResponse(fw.CodeFailedSerialization, "could not validate provided json")
			fw.JError(w, resp, h.StatusBadRequest)
			return
		}

		err := admin.CreateDataSource([]*graphx.DataSource{&v})
		if err != nil {
			logger.Error().Msgf("failed to create datasource %v: %v", v.Name, err)
			resp := fw.NewResponse(fw.CodeCreateFail, "creation failed")
			fw.JError(w, resp, h.StatusInternalServerError)
			return
		}
		logger.Debug().Msgf("successfully created datasource: %v", v.Name)

		resp := fw.NewResponse(fw.CodeSuccess, "DataSource stored")
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			logger.Error().Msgf("failed to serialize response: %v", err)
			w.WriteHeader(h.StatusInternalServerError)
			return
		}

		return
	}
}
