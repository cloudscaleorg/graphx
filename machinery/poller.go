package machinery

import (
	"context"
	"log"
	"time"

	"github.com/cloudscaleorg/graphx"
)

// Poller is reusable machinery which calls a Querier's query method at a specific interval.
// this is useful for backend that do not provide a native streaming api such as prometheus.
type Poller struct {
	// ID representing the unique session with a client
	ID string
	// the instance of a querier implementation used to query the database
	Q graphx.Querier
	// the interval in which we call Query() on the querier
	PollInterval time.Duration
}

// NewPoller is a contructor for a poller.
func NewPoller(id string, q graphx.Querier, pollInterval time.Duration) *Poller {
	return &Poller{
		ID:           id,
		Q:            q,
		PollInterval: pollInterval,
	}
}

// Poll is intended to be ran as a go routine and will call it's Querier query method.
func (p *Poller) Poll(ctx context.Context) {
	t := time.NewTicker(p.PollInterval)
	defer t.Stop()

	log.Printf("poller id %s: beginning polling at %v", p.ID, p.PollInterval)
	for {
		select {
		case <-ctx.Done():
			log.Printf("poller id %s: context cancled. polling stopped", p.ID)
			return
		case <-t.C:
			startTS := time.Now()
			p.Q.Query(ctx)
			endTS := time.Now().Sub(startTS)
			log.Printf("poller id %s: all queries to datastore took %v", p.ID, endTS)
		}
	}
}
