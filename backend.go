package graphx

// Backend represents an implemented DataSource backend.
type Backend interface {
	// a unique name for this backend
	Name() string
	// a constructor for a Querier
	Querier(metrics []*ChartMetric) Querier
}
