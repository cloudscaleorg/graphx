package graphx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ldelossa/jsonerr"
	validator "gopkg.in/go-playground/validator.v9"
)

const (
	ValidationError      = "could not validate your charts descriptor. chart_names and names keys are required and must contain more the one item"
	MetricsStreamErrCode = "graphx.stream_handler"
)

func StreamHandler(v *validator.Validate, cs ChartStore, sf StreamerFactory, ws websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only support posts
		if r.Method != http.MethodGet {
			log.Printf("methd not allowed")
			resp := jsonerr.NewResponse("", MetricsStreamErrCode, "method not allowed")
			jsonerr.Error(w, resp, http.StatusMethodNotAllowed)
			return
		}

		// create context for this streaming session
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// upgrade to web socket
		wsConn, err := ws.Upgrade(w, r, nil)
		defer wsConn.Close()

		if err != nil {
			log.Printf("failed to upgrade to websocket: %v", err)
			resp := jsonerr.NewResponse("", MetricsStreamErrCode, "failed to upgrade to websocket")
			jsonerr.Error(w, resp, http.StatusBadRequest)
			return
		}
		log.Printf("successfully upgraded to websocket")

		// TODO: handle timeouts
		// set initial deadline see: https://github.com/golang/go/blob/master/src/net/net.go#L149

		// wait for charts descriptor
		var cd ChartsDescriptor
		for {
			err := wsConn.ReadJSON(&cd)
			if err != nil {
				log.Printf("received error waiting for chart descriptor: %v", err)
				return
			}
			break
		}
		id := fmt.Sprintf("%s.%v", uuid.New().String(), cd.Names)
		log.Printf("id: %v received chart descriptor: %v", id, cd)

		// validate struct
		err = v.StructCtx(ctx, cd)
		if err != nil {
			log.Printf("id %s: struct validation error: %v", id, err)
			return
		}

		// do not allow polls of lower then a second
		if time.Duration(cd.PollInterval) < 1*time.Second {
			log.Printf("id %s: requested poll interval of less then 1 second", id)
			return
		}

		// receive configured charts from chart store
		charts, err := cs.GetByNames(cd.ChartNames)
		if err != nil {
			log.Printf("id %s: failed to query chart store: %v", id, err)
			return
		}

		// create streamer from our streamer factory
		_ = sf.NewStreamer(ctx, id, charts, time.Duration(cd.PollInterval))
		if err != nil {
			log.Printf("id %s: failed to instantiate query streamer: %v", id, err)
			return
		}

	}
}
