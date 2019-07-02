package prometheus

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/cloudscaleorg/graphx"
	promapi "github.com/prometheus/client_golang/api/prometheus/v1"
	prommodels "github.com/prometheus/common/model"
)

type QuerierOpts struct {
	ID           string
	Client       promapi.API
	ChartMetrics []graphx.ChartMetric
	MChan        chan *graphx.Metric
	EChan        chan error
}

type querier struct {
	QuerierOpts
}

// NewQuerier creates a prometheus Querier.
func NewQuerier(opts QuerierOpts) graphx.Querier {
	q := &querier{
		QuerierOpts: opts,
	}

	return q
}

// Query is the public method implementing the graphx.Querier interface. this method blocks
// until all concurrent queries are completed and have streamed their metrics to the provided channel
func (q *querier) Query(ctx context.Context) {
	var wg sync.WaitGroup

	// check context
	select {
	case <-ctx.Done():
		log.Printf("session id %s: context closed before concurrent queries", q.ID)
		return
	default:
	}

	// derive context with timeout
	ctxTO, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for _, chartMetric := range q.ChartMetrics {
		wg.Add(1)
		go q.query(ctxTO, chartMetric.Name, chartMetric.Query, &wg)
	}

	wg.Wait()
}

// query is a private method meant to be ran as a go routine. handles the logic for querying prometheus given
// a chart and a query and streams the results to the internal metrics channel
func (q *querier) query(ctx context.Context, chart string, query string, wg *sync.WaitGroup) {
	defer wg.Done()

	// check context
	select {
	case <-ctx.Done():
		log.Printf("session id %s: context closed before query", q.ID)
		return
	default:
	}

	// issue query
	value, _, err := q.Client.Query(ctx, string(query), time.Now())
	if err != nil {
		log.Printf("session id %s: failed to query prometheus. ERROR: %v QUERY: %v", q.ID, err, string(query))
		return
	}

	// type assert return value to vector
	var vector prommodels.Vector
	var ok bool
	if vector, ok = value.(prommodels.Vector); !ok {
		log.Printf("received unknown type from vector request")
		return
	}

	// stream metrics to channel
	for _, sample := range vector {
		m := sampleToMetric(chart, sample)
		select {
		case <-ctx.Done():
			log.Printf("session id %s: context closed while streaming results", q.ID)
			return
		case q.MChan <- m:
		default:
			log.Printf("session id %s: unable to deliver metrics to channel", q.ID)
		}
	}

}
