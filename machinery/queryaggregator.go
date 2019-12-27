package machinery

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/cloudscaleorg/graphx"
)

// QueryAggregator acts as a fan-in for multiple Queriers
//
//
// QueryAggregator implements the graphx.Streamer interface
// Each provided Querier will have it's Query method called at queryInterval interval
// The queried metrics are returned to the internal mChan and maybe retrieved via the Recv method
type QueryAggregator struct {
	queriers      []graphx.Querier
	mChan         chan *graphx.Metric
	eChan         chan error
	queryInterval time.Duration
}

func NewQueryAggregator(queriers []graphx.Querier, interval time.Duration) *QueryAggregator {
	mChan := make(chan *graphx.Metric, 1024)
	eChan := make(chan error, 1024)
	return &QueryAggregator{
		queriers:      queriers,
		mChan:         mChan,
		eChan:         eChan,
		queryInterval: interval,
	}
}

func (qa *QueryAggregator) Recv(ctx context.Context) (*graphx.Metric, error) {
	select {
	case m := <-qa.mChan:
		return m, nil
	case err := <-qa.eChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (qa *QueryAggregator) Start(ctx context.Context) {
	go qa.start(ctx)
}

func (qa *QueryAggregator) start(ctx context.Context) {
	qa.query(ctx)
	for {
		ticker := time.NewTicker(qa.queryInterval)
		select {
		case <-ticker.C:
			qa.query(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (qa *QueryAggregator) query(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for _, q := range qa.queriers {
		wg.Add(1)
		qq := q
		go func() {
			defer wg.Done()
			metrics, err := qq.Query(ctx)
			if err != nil {
				select {
				case qa.eChan <- err:
				default:
					log.Printf("could not return error %v", err)
					return
				}
			}
			for _, m := range metrics {
				select {
				case qa.mChan <- m:
				default:
					log.Printf("could not write metric %v", m)
				}
			}
		}()
	}
	wg.Done()
}
