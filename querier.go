package graphx

import (
	"context"
)

// Querier will run one or more ChartMetric queries for the implemented Backend and return Metric structs to the caller
type Querier interface {
	Query(ctx context.Context) ([]*Metric, error)
}
