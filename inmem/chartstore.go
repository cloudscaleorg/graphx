package inmem

import (
	"sync"

	"github.com/cloudscaleorg/graphx"
)

type chartStore struct {
	mu *sync.RWMutex
	m  map[string]*graphx.Chart
}

func NewChartStore() graphx.ChartStore {
	return &chartStore{
		mu: &sync.RWMutex{},
		m:  make(map[string]*graphx.Chart),
	}
}

func (cs *chartStore) Get() ([]*graphx.Chart, error) {
	// allocate
	a := make([]*graphx.Chart, 0)

	// copy
	cs.mu.RLock()
	for _, v := range cs.m {
		a = append(a, v)
	}
	cs.mu.RUnlock()

	return a, nil
}

func (cs *chartStore) GetByNames(chartNames []string) ([]*graphx.Chart, error) {
	// allocate
	a := make([]*graphx.Chart, 0)

	// retrieve
	for _, chartName := range chartNames {
		cs.mu.RLock()
		a = append(a, cs.m[chartName])
		cs.mu.RUnlock()
	}

	return a, nil
}

func (cs *chartStore) Store(charts []*graphx.Chart) error {
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
