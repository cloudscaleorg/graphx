package http

import (
	"encoding/json"
	h "net/http"

	"github.com/cloudscaleorg/graphx"
	"github.com/cloudscaleorg/graphx/admin"
	fw "github.com/ldelossa/goframework/http"
	"github.com/rs/zerolog/log"
)

func UpdateDataSource(a admin.DataSource) h.HandlerFunc {
	logger := log.With().Str("component", "UpdateDataSourceHandler").Logger()
	return func(w h.ResponseWriter, r *h.Request) {
		if r.Method != h.MethodPut {
			resp := fw.NewResponse(fw.CodeMethodNotImplemented, "endpoint only supports PUT")
			fw.JError(w, resp, h.StatusNotImplemented)
			return
		}

		var v graphx.DataSource
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			logger.Error().Msgf("could not validate provided json: %v", err)
			resp := fw.NewResponse(fw.CodeFailedSerialization, "could not validate provided json")
			fw.JError(w, resp, h.StatusBadRequest)
			return
		}

		err = a.UpdateDataSource(&v)
		if err != nil {
			switch {
			case (err == admin.ErrNotFound):
				logger.Error().Msg("resource being updated not found")
				resp := fw.NewResponse(fw.CodeNotFound, "resource being updated not found")
				fw.JError(w, resp, h.StatusNotFound)
				return
			case (err == admin.ErrStore{}):
				logger.Error().Msgf("storage error: %v", err)
				resp := fw.NewResponse(fw.CodeInternalServerError, "an internal error occured")
				fw.JError(w, resp, h.StatusNotFound)
				return
			}
		}

		err = json.NewEncoder(w).Encode(&v)
		if err != nil {
			w.WriteHeader(h.StatusInternalServerError)
			return
		}

		logger.Debug().Msgf("succesfully updated datasource: %v", v.Name)
		return
	}
}
