package machinery

import (
	"context"
	"time"

	"github.com/cloudscaleorg/graphx"
	"github.com/cloudscaleorg/graphx/registry"
)

// Aggregator implements the Streamer interface.
// Aggregator is instantiated with a list of Charts, merges those charts
// and then starts the necessary Queries. each Querier is provided a
// metrics and an error channel to communicate with the Aggregator.
type aggregator struct {
	AggregatorOpts
	// an id representing this streaming session
	id string
	// the metrics channel Queriers will deliver metrics on
	mChan chan *graphx.Metric
	// the error channel Queriers will deliver errors on
	eChan chan error
	// a registry for constructing Queriers
	reg registry.Querier
}

// AggregatorOpts are the options for an aggregator
type AggregatorOpts struct {
	PollInterval time.Duration
	Charts       []*graphx.Chart
	ChartMetrics map[string][]*graphx.ChartMetric
	Registry     registry.Querier
}

// aggregate will create multiple queriers.
// datasources have the conn string necessary to create a client
// chart metrics have the datasource.
// aggregator gets list of chart metrics, looksups data sources
// creates a querier per data source, datasource has "type".

// NewAggregator creates an aggregator Streamer. make sure to cancel ctx
// to not leak go routines.
func NewAggregator(ctx context.Context, id string, opts AggregatorOpts) graphx.Streamer {
	// TODO: determine best size for buffered channel.
	mChan := make(chan *graphx.Metric, 1024)
	eChan := make(chan error, 1024)

	// chartMetrics := graphx.DatasourceTranspose(opts.Charts)

	// for datasource, chartMetrics := range chartMetrics {
	// 	switch datasource {
	// 	case prometheus.Datasource:
	// 		pOpts := prometheus.QuerierOpts{
	// 			ID:           id,
	// 			Client:       opts.PromClient,
	// 			ChartMetrics: chartMetrics,
	// 			MChan:        mChan,
	// 			EChan:        eChan,
	// 		}

	// 		// create a prometheus querier and a poller and launch
	// 		pq := prometheus.NewQuerier(pOpts)
	// 		poller := NewPoller(id, pq, opts.PollInterval)
	// 		go poller.Poll(ctx)
	// 	case "influxdb":
	// 	}
	// }

	return &aggregator{
		AggregatorOpts: opts,
		id:             id,
		mChan:          mChan,
		eChan:          eChan,
	}
}

func (a *aggregator) Recv() (*graphx.Metric, error) {
	select {
	case m := <-a.mChan:
		return m, nil
	case e := <-a.eChan:
		return nil, e
	}
}
