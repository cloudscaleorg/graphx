package graphx

import "context"

// Backend represents an implemented DataSource backend.
//
// Backends are responsible for constructing various features of a backend such as Queriers
type Backend interface {
	// a unique name for this backend
	Name() string
	// a constructor for a Querier
	Querier(metrics []*ChartMetric) Querier
}

// Querier performs a query for each configured ChartMetric and returns the aggregation
type Querier interface {
	Query(ctx context.Context) ([]*Metric, error)
}
