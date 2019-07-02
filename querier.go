package graphx

import (
	"context"
)

// Querier is an interface to abstract the data backend we retrieve metrics from.
// implementations can control how they query the backend data source
type Querier interface {
	// Query the results of a query to the provided channel
	Query(ctx context.Context)
}
