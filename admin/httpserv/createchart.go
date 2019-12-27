package httpserv

import (
	"encoding/json"
	h "net/http"

	"github.com/cloudscaleorg/graphx"
	"github.com/cloudscaleorg/graphx/admin"
	fw "github.com/ldelossa/goframework/http"
	"github.com/rs/zerolog/log"
)

func CreateChart(admin *admin.Admin) h.HandlerFunc {
	logger := log.With().Str("component", "CreateChartHandler").Logger()
	return func(w h.ResponseWriter, r *h.Request) {
		if r.Method != h.MethodPost {
			resp := fw.NewResponse(fw.CodeMethodNotImplemented, "endpoint only supports POST")
			fw.JError(w, resp, h.StatusNotImplemented)
			return
		}
		var v graphx.Chart
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			logger.Error().Msgf("failed to deserialize datasource: %v", err)
			resp := fw.NewResponse(fw.CodeFailedSerialization, "could not validate provided json")
			fw.JError(w, resp, h.StatusBadRequest)
			return
		}
		err := admin.CreateChart([]*graphx.Chart{&v})
		if err != nil {
			logger.Error().Msgf("failed to create chart: %v", err)
			resp := fw.NewResponse(fw.CodeCreateFail, err.Error())
			fw.JError(w, resp, h.StatusBadRequest)
			return
		}
		resp := fw.NewResponse(fw.CodeSuccess, "DataSource stored")
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			logger.Error().Msgf("failed to serialize response: %v", err)
			w.WriteHeader(h.StatusInternalServerError)
			return
		}
		logger.Debug().Msgf("successfully created chart: %v", v.Name)
		return
	}
}
