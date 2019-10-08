package machinery

import (
	"context"
	"log"
	"time"

	"github.com/cloudscaleorg/graphx"
)

// QueryAggregator
type QueryAggregator struct {
	queriers      []graphx.Querier
	mChan         chan *graphx.Metric
	eChan         chan error
	queryInterval time.Duration
}

func (qa *QueryAggregator) Start(ctx context.Context) error {

	for _, q := range qa.queriers {
		go func(q graphx.Querier) {
			metrics, err := q.Query(ctx)
			if err != nil {
				select {
				case qa.eChan <- err:
					return
				default:
					log.Printf("could not send error")
					return
				}
			}

			for _, metric := range metrics {
				qa.mChan <- metric
			}
		}(q)
	}
}
