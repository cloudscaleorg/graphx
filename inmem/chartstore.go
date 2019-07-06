package inmem

import (
	"sync"

	"github.com/cloudscaleorg/graphx"
)

type chartStore struct {
	mu *sync.RWMutex
	m  map[string]graphx.Chart
}

func NewChartStore() graphx.ChartStore {
	return &chartStore{
		mu: &sync.RWMutex{},
		m:  make(map[string]graphx.Chart),
	}
}

func (cs *chartStore) Get() (map[string]graphx.Chart, error) {
	// allocate
	m := make(map[string]graphx.Chart)

	// copy
	cs.mu.RLock()
	for k, v := range cs.m {
		m[k] = v
	}
	cs.mu.RUnlock()

	return m, nil
}

func (cs *chartStore) GetByNames(chartNames []string) (map[string]graphx.Chart, error) {
	// allocate
	m := make(map[string]graphx.Chart)

	// retrieve
	for _, chartName := range chartNames {
		cs.mu.RLock()
		m[chartName] = cs.m[chartName]
		cs.mu.RUnlock()
	}

	return m, nil
}

func (cs *chartStore) Store(charts []graphx.Chart) error {
	for _, chart := range charts {
		cs.mu.Lock()
		cs.m[chart.Name] = chart
		cs.mu.Unlock()
	}

	return nil
}

func (cs *chartStore) RemoveByNames(chartNames []string) error {
	for _, chartName := range chartNames {
		cs.mu.Lock()
		delete(cs.m, chartName)
		cs.mu.Unlock()
	}

	return nil
}
