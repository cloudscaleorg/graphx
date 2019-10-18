package prometheus

import "github.com/cloudscaleorg/graphx"

const (
	Prometheus = "prometheus"
)

// Backend implements the graphx.Backend interface for Prometheus.
//
// On Backend construction a Prometheus client will be configured to
// connect to the graphx.DataSource.ConnString.
//
// After successful backend construction a "feature" of the backend such as the Querier
// constructor maybe called. The returned Querier will send all ChartMetric.Query strings
// to the configured DataSource.
type Backend struct {
}

// implements the registry.BackendFactory function signature
func NewBackend(ds *graphx.DataSource) (graphx.Backend, error) {
	panic("not implemented")
}

// Name is the name of this Backend
func (b *Backend) Name() string {
	panic("not implemented")
}

// Querier returns a Prometheus Querier.
func (b *Backend) Querier(metrics []*graphx.ChartMetric) graphx.Querier {
	panic("not implemented")
}
