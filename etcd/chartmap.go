package etcd

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudscaleorg/events"
	"github.com/cloudscaleorg/graphx"
	"github.com/rs/zerolog"
	v3 "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

var cReduce = func(m *ChartMap) events.ReduceFunc {
	return func(e *v3.Event, snapshot bool) {
		if snapshot {
			m.Reset()
		}
		switch e.Type {
		case mvccpb.PUT:
			var v graphx.Chart
			err := v.FromJSON(e.Kv.Value)
			if err != nil {
				m.logger.Error().Msgf("failed to deserialize key %v: %v", string(e.Kv.Key), err)
				return
			}
			m.Store([]*graphx.Chart{&v})
			m.logger.Debug().Msgf("added datasource to store: %v", v.Name)
		case mvccpb.DELETE:
			var v graphx.Chart
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

type ChartMap struct {
	mu       *sync.RWMutex
	m        map[string]*graphx.Chart
	etcd     *v3.Client
	listener *events.Listener
	logger   zerolog.Logger
}

func NewChartMap(ctx context.Context, client *v3.Client) (*ChartMap, error) {
	m := &ChartMap{
		mu:   &sync.RWMutex{},
		m:    map[string]*graphx.Chart{},
		etcd: client,
	}
	l, err := events.NewListener(&events.Opts{
		Prefix: cPrefix,
		Client: client,
		F:      cReduce(m),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create events listener: %v", err)
	}
	m.listener = l
	l.Listen(ctx)
	l.Ready(ctx)
	return m, nil
}

func (m *ChartMap) Get(names []string) ([]*graphx.Chart, []string) {
	out := []*graphx.Chart{}
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

func (m *ChartMap) Remove(names []string) {
	m.mu.Lock()
	for _, name := range names {
		delete(m.m, name)
	}
	m.mu.Unlock()
}

func (m *ChartMap) Store(charts []*graphx.Chart) {
	m.mu.Lock()
	for _, chart := range charts {
		m.m[chart.Name] = chart
	}
	m.mu.Unlock()
}

func (m *ChartMap) Reset() {
	m.mu.Lock()
	m.m = map[string]*graphx.Chart{}
	m.mu.Unlock()
}
