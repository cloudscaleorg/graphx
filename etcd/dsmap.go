package etcd

import (
	"sync"

	"github.com/cloudscaleorg/graphx"
)

type dsmap struct {
	mu *sync.RWMutex
	m  map[string]*graphx.DataSource
}

func NewDSMap() *dsmap {
	return &dsmap{
		mu: &sync.RWMutex{},
		m:  map[string]*graphx.DataSource{},
	}
}

func (m *dsmap) get(names []string) ([]*graphx.DataSource, []string) {
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

func (m *dsmap) remove(names []string) {
	m.mu.Lock()
	for _, name := range names {
		delete(m.m, name)
	}
	m.mu.Unlock()
}

func (m *dsmap) store(charts []*graphx.DataSource) {
	m.mu.Lock()
	for _, chart := range charts {
		m.m[chart.Name] = chart
	}
	m.mu.Unlock()
}

func (m *dsmap) reset() {
	m.mu.Lock()
	m.m = map[string]*graphx.DataSource{}
	m.mu.Unlock()
}
