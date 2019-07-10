package machinery

import (
	"context"
	"time"

	"github.com/cloudscaleorg/graphx"
	promapi "github.com/prometheus/client_golang/api/prometheus/v1"
)

// aggregatorFactory holds any constant runtime depedencies for an aggregator streamer.
type aggregatorFactory struct {
	// prometheus client
	pc promapi.API
	// influx client
	// opentsd client
	// ...
}

func NewAggregatorFactory(promClient promapi.API) graphx.StreamerFactory {
	return &aggregatorFactory{
		pc: promClient,
	}
}

func (af *aggregatorFactory) NewStreamer(ctx context.Context, id string, charts []*graphx.Chart, pollInterval time.Duration) graphx.Streamer {
	opts := AggregatorOpts{
		PollInterval: pollInterval,
		Charts:       charts,
		PromClient:   af.pc,
	}

	streamer := NewAggregator(ctx, id, opts)

	return streamer
}
