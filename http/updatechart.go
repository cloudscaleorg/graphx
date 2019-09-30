package http

import (
	"encoding/json"
	"fmt"
	h "net/http"

	"github.com/cloudscaleorg/graphx"
	"github.com/cloudscaleorg/graphx/admin"
	fw "github.com/ldelossa/goframework/http"
	"github.com/rs/zerolog/log"
)

func UpdateChart(a admin.Chart) h.HandlerFunc {
	logger := log.With().Str("component", "UpdateChartHandler").Logger()
	return func(w h.ResponseWriter, r *h.Request) {
		if r.Method != h.MethodPut {
			resp := fw.NewResponse(fw.CodeMethodNotImplemented, "endpoint only supports PUT")
			fw.JError(w, resp, h.StatusNotImplemented)
			return
		}

		var v graphx.Chart
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			logger.Error().Msgf("failed to deserialize datasource: %v", err)
			resp := fw.NewResponse(fw.CodeFailedSerialization, "could not validate provided json")
			fw.JError(w, resp, h.StatusBadRequest)
			return
		}

		err = a.UpdateChart(&v)
		switch e := err.(type) {
		case admin.ErrNotFound:
			logger.Error().Msgf("resource being updated not found: %v", e)
			resp := fw.NewResponse(fw.CodeNotFound, "resource being updated not found")
			fw.JError(w, resp, h.StatusNotFound)
			return
		case admin.ErrStore:
			logger.Error().Msgf("storage error: %v", err)
			resp := fw.NewResponse(fw.CodeInternalServerError, "an internal error occured")
			fw.JError(w, resp, h.StatusNotFound)
			return
		case admin.ErrMissingDataSources:
			log.Printf("got here")
			logger.Error().Msgf("missing datasource for %v: %v", v.Name, e)
			s := fmt.Sprintf("missing datasources: %v", err.Error())
			resp := fw.NewResponse(fw.CodeNotFound, s)
			fw.JError(w, resp, h.StatusNotFound)
			return
		default:
			logger.Error().Msgf("failed to update chart %v: %v", v.Name, e)
			resp := fw.NewResponse(fw.CodeNotFound, "")
			fw.JError(w, resp, h.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(&v)
		if err != nil {
			logger.Error().Msgf("failed to deserialize datasource: %v", err)
			w.WriteHeader(h.StatusInternalServerError)
			return
		}

		logger.Debug().Msgf("successfully updated chart: %v", v.Name)
		return
	}
}
