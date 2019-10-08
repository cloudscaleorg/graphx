package prometheus

import "github.com/cloudscaleorg/graphx"

// Backend implements the graphx.Backend interface for Prometheus
type Backend struct {
}

// implements the registry.BackendFactory function signature
func NewBackend(ds *graphx.DataSource) (*graphx.Backend, error) {
	panic("not implemented")
}

// Name is the name of this Backend
func (b *Backend) Name() string {
	panic("not implemented")
}

// Querier returns a Prometheus Querier implementation.
func (b *Backend) Querier(metrics []*graphx.ChartMetric) graphx.Querier {
	panic("not implemented")
}
