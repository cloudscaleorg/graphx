package etcd

import (
	"context"
	"fmt"

	"github.com/cloudscaleorg/events"
	"github.com/cloudscaleorg/graphx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	v3 "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

var cReduce = func(store *ChartStore) events.ReduceFunc {
	return func(e *v3.Event, snapshot bool) {
		if snapshot {
			store.reset()
		}

		switch e.Type {
		case mvccpb.PUT:
			var v graphx.Chart
			err := v.FromJSON(e.Kv.Value)
			log.Printf("after return: %v", v)
			if err != nil {
				store.logger.Error().Msgf("failed to deserialize key %v: %v", string(e.Kv.Key), err)
				return
			}
			store.store([]*graphx.Chart{&v})
			store.logger.Debug().Msgf("added datasource to store: %v", v.Name)
		case mvccpb.DELETE:
			var v graphx.Chart
			err := v.FromJSON(e.PrevKv.Value)
			if err != nil {
				store.logger.Error().Msgf("failed to deserialize previous key %v: %v", string(e.Kv.Key), err)
				return
			}
			store.remove([]string{v.Name})
			store.logger.Debug().Msgf("removed datasource from store: %v", v.Name)
		}
	}
}

// ChatStore implement the graphx.ChatStore interface
type ChartStore struct {
	*chartmap
	etcd     *v3.Client
	listener *events.Listener
	logger   zerolog.Logger
}

func NewChartStore(ctx context.Context, client *v3.Client) (graphx.ChartStore, error) {
	store := &ChartStore{
		chartmap: NewChartMap(),
		etcd:     client,
		logger:   log.With().Str("component", "chartstore").Logger(),
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
	out, missing := s.get(nil)
	if len(missing) > 0 {
		return nil, &ErrNotFound{
			missing: missing,
		}
	}

	return out, nil
}

func (s *ChartStore) GetByNames(names []string) ([]*graphx.Chart, []string, error) {
	out, missing := s.get(names)
	return out, missing, nil
}

func (s *ChartStore) Store(charts []*graphx.Chart) error {
	if err := s.checkListener(); err != nil {
		return err
	}
	err := putCharts(context.Background(), s.etcd, charts)
	return err
}

func (s *ChartStore) RemoveByNames(names []string) error {
	if err := s.checkListener(); err != nil {
		return err
	}
	err := delCharts(context.Background(), s.etcd, names)
	return err
}

func (s *ChartStore) checkListener() error {
	ctx, cancel := context.WithTimeout(context.Background(), readyTO)
	defer cancel()

	err := s.listener.Ready(ctx)
	return err
}
