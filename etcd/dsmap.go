package etcd

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudscaleorg/events"
	"github.com/cloudscaleorg/graphx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	v3 "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// dsReduce is an events.ReduceFunc which updates a local
// DSMap on receipt of events occuring on the DataSource etcd prefix
var dsReduce = func(m *DSMap) events.ReduceFunc {
	return func(e *v3.Event, snapshot bool) {
		// a new consistent state will follow. dump the current map and get up to date
		if snapshot {
			m.Reset()
		}
		switch e.Type {
		case mvccpb.PUT:
			var v graphx.DataSource
			err := v.FromJSON(e.Kv.Value)
			if err != nil {
				m.logger.Error().Msgf("failed to deserialize key %v: %v", string(e.Kv.Key), err)
				return
			}
			m.Store([]*graphx.DataSource{&v})
			m.logger.Debug().Msgf("added datasource to store: %v", v.Name)
		case mvccpb.DELETE:
			var v graphx.DataSource
			err := v.FromJSON(e.PrevKv.Value)
			if err != nil {
				m.logger.Error().Msgf("failed to deserialize previous key %v: %v", string(e.Kv.Key), err)
				return
			}
			m.Remove([]string{v.Name})
			m.logger.Debug().Msgf("removed datasource from store: %v", v.Name)
		}
	}
}

// DSMap keeps a local map of DataSources existing in a GraphX cluster.
//
// DSMap's state is kept in sync with etcd in an eventually consistent fashion.
type DSMap struct {
	mu       *sync.RWMutex
	m        map[string]*graphx.DataSource
	etcd     *v3.Client
	listener *events.Listener
	logger   zerolog.Logger
}

// NewDSMap is a constructor for a DSMap
func NewDSMap(ctx context.Context, client *v3.Client) (*DSMap, error) {
	m := &DSMap{
		mu:     &sync.RWMutex{},
		m:      map[string]*graphx.DataSource{},
		etcd:   client,
		logger: log.With().Str("component", "datasourcestore").Logger(),
	}
	l, err := events.NewListener(&events.Opts{
		Prefix: dsPrefix,
		Client: client,
		F:      dsReduce(m),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create events listener: %v", err)
	}
	m.listener = l
	l.Listen(ctx)
	l.Ready(ctx)
	return m, nil
}

// Get retrieves graphx.DataSource objects.
//
// If a nil array if provided all DataSources are returned.
// If an array of strings are provided only matching DataSources will be returned.
// If the array contains names that do not exist in the state these names will be returned
// as the second argument indicating the missing DataSources.
func (m *DSMap) Get(names []string) ([]*graphx.DataSource, []string) {
	out := []*graphx.DataSource{}
	missing := []string{}
	if names == nil {
		m.mu.RLock()
		for _, v := range m.m {
			out = append(out, v)
		}
		m.mu.RUnlock()
		return out, missing
	}
	m.mu.RLock()
	for _, name := range names {
		v, ok := m.m[name]
		if !ok {
			missing = append(missing, name)
		} else {
			out = append(out, v)
		}
	}
	m.mu.RUnlock()
	return out, missing
}

// Remove removes any DataSources matching the provided
// names
//
// Names which do not exist are ignored.
func (m *DSMap) Remove(names []string) {
	m.mu.Lock()
	for _, name := range names {
		delete(m.m, name)
	}
	m.mu.Unlock()
}

// Store persists a slice of DataSources
func (m *DSMap) Store(charts []*graphx.DataSource) {
	m.mu.Lock()
	for _, chart := range charts {
		m.m[chart.Name] = chart
	}
	m.mu.Unlock()
}

// Reset discards the current DSMap state
func (m *DSMap) Reset() {
	m.mu.Lock()
	m.m = map[string]*graphx.DataSource{}
	m.mu.Unlock()
}
