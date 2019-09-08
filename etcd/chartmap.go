package etcd

import (
	"log"
	"sync"

	"github.com/cloudscaleorg/graphx"
)

type chartmap struct {
	mu *sync.RWMutex
	m  map[string]*graphx.Chart
}

func NewChartMap() *chartmap {
	return &chartmap{
		mu: &sync.RWMutex{},
		m:  map[string]*graphx.Chart{},
	}
}

func (m *chartmap) get(names []string) ([]*graphx.Chart, []string) {
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

func (m *chartmap) remove(names []string) {
	m.mu.Lock()
	for _, name := range names {
		delete(m.m, name)
	}
	m.mu.Unlock()
}

func (m *chartmap) store(charts []*graphx.Chart) {
	m.mu.Lock()
	for _, chart := range charts {
		log.Printf("%v", chart)
		m.m[chart.Name] = chart
	}
	m.mu.Unlock()
}

func (m *chartmap) reset() {
	m.mu.Lock()
	m.m = map[string]*graphx.Chart{}
	m.mu.Unlock()
}
