package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cloudscaleorg/events"
	"github.com/cloudscaleorg/graphx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	v3 "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

const (
	cPrefix         = "/graphx/charts"
	cPrefixTemplate = "/graphx/charts/%s"
	readyTO         = 5 * time.Second
)

// dsReduce closes over a DSStore and returns a function for
// reducing events
var cReduce = func(store *ChartStore) events.ReduceFunc {
	return func(e *v3.Event) {
		switch e.Type {
		case mvccpb.PUT:
			var v *graphx.Chart
			if err := json.Unmarshal(e.Kv.Value, &v); err != nil {
				store.logger.Error().Msgf("failed to unmarshal event for key: %v", string(e.Kv.Key))
				return
			}
			store.mu.Lock()
			store.m[v.Name] = v
			store.mu.Unlock()
			store.logger.Debug().Msgf("added datasource to store: %v", v.Name)
		case mvccpb.DELETE:
			var v *graphx.Chart
			if err := json.Unmarshal(e.PrevKv.Value, &v); err != nil {
				store.logger.Error().Msgf("failed to unmarshal prev event for key: %v", string(e.Kv.Key))
				return
			}
			store.mu.Lock()
			delete(store.m, v.Name)
			store.mu.Unlock()
			store.logger.Debug().Msgf("removed datasource from store: %v", v.Name)
		}
	}
}

// dsStore implement the graphx.DataSourceStore
type ChartStore struct {
	etcd     *v3.Client
	mu       *sync.RWMutex
	m        map[string]*graphx.Chart
	listener *events.Listener
	logger   zerolog.Logger
}

func NewChartStore(ctx context.Context, client *v3.Client) (graphx.ChartStore, error) {
	store := &ChartStore{
		etcd:   client,
		mu:     &sync.RWMutex{},
		m:      map[string]*graphx.Chart{},
		logger: log.With().Str("component", "chartstore").Logger(),
	}

	l, err := events.NewListener(&events.Opts{
		Prefix: cPrefix,
		Client: client,
		F:      cReduce(store),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create events listener: %v", err)
	}
	store.listener = l

	l.Listen(ctx)
	l.Ready(ctx)

	return store, nil
}

func (s *ChartStore) Get() ([]*graphx.Chart, error) {
	if err := s.checkListener(); err != nil {
		return nil, err
	}

	d := []*graphx.Chart{}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, v := range s.m {
		d = append(d, v)
	}
	return d, nil
}

func (s *ChartStore) GetByNames(names []string) ([]*graphx.Chart, error) {
	if err := s.checkListener(); err != nil {
		return nil, err
	}

	d := []*graphx.Chart{}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, name := range names {
		if v, ok := s.m[name]; ok {
			d = append(d, v)
		}
	}

	return d, nil
}
func (s *ChartStore) Store(charts []*graphx.Chart) error {
	if err := s.checkListener(); err != nil {
		return err
	}

	for _, chart := range charts {
		b, err := json.Marshal(chart)
		if err != nil {
			return fmt.Errorf("failed to serialize datastore: %v", err)
		}

		key := fmt.Sprintf(cPrefixTemplate, chart.Name)
		val := string(b)
		_, err = s.etcd.Put(context.Background(), key, val)
		if err != nil {
			return fmt.Errorf("failed to add ds %v to store", chart.Name)
		}
	}

	return nil
}
func (s *ChartStore) RemoveByNames(names []string) error {
	if err := s.checkListener(); err != nil {
		return err
	}

	for _, name := range names {
		_, err := s.etcd.Delete(context.Background(), name)
		if err != nil {
			return fmt.Errorf("failed to remove ds %v to store", name)
		}
	}

	return nil
}

func (s *ChartStore) checkListener() error {
	ctx, cancel := context.WithTimeout(context.Background(), readyTO)
	defer cancel()

	err := s.listener.Ready(ctx)
	return err
}
