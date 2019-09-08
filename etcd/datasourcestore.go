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

var dsReduce = func(store *DSStore) events.ReduceFunc {
	return func(e *v3.Event, snapshot bool) {
		// a new consistent state will follow. dump the current map and get up to date
		if snapshot {
			store.reset()
		}

		switch e.Type {
		case mvccpb.PUT:
			var v graphx.DataSource
			err := v.FromJSON(e.Kv.Value)
			if err != nil {
				store.logger.Error().Msgf("failed to deserialize key %v: %v", string(e.Kv.Key), err)
				return
			}
			store.store([]*graphx.DataSource{&v})
			store.logger.Debug().Msgf("added datasource to store: %v", v.Name)
		case mvccpb.DELETE:
			var v *graphx.DataSource
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

// dsStore implement the graphx.DataSourceStore
type DSStore struct {
	*dsmap
	etcd     *v3.Client
	listener *events.Listener
	logger   zerolog.Logger
}

func NewDSStore(ctx context.Context, client *v3.Client) (graphx.DataSourceStore, error) {
	store := &DSStore{
		dsmap:  NewDSMap(),
		etcd:   client,
		logger: log.With().Str("component", "datasourcestore").Logger(),
	}

	l, err := events.NewListener(&events.Opts{
		Prefix: dsPrefix,
		Client: client,
		F:      dsReduce(store),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create events listener: %v", err)
	}
	store.listener = l

	l.Listen(ctx)
	l.Ready(ctx)

	return store, nil
}

func (s *DSStore) Get() ([]*graphx.DataSource, error) {
	out, missing := s.get(nil)
	if len(missing) > 0 {
		return nil, &ErrNotFound{
			missing: missing,
		}
	}

	return out, nil
}

func (s *DSStore) GetByNames(names []string) ([]*graphx.DataSource, error) {
	out, missing := s.get(names)
	if len(missing) > 0 {
		return nil, &ErrNotFound{
			missing: missing,
		}
	}
	return out, nil
}

func (s *DSStore) Store(sources []*graphx.DataSource) error {
	if err := s.checkListener(); err != nil {
		return err
	}
	err := putDS(context.Background(), s.etcd, sources)
	return err
}

func (s *DSStore) RemoveByNames(names []string) error {
	if err := s.checkListener(); err != nil {
		return err
	}
	err := delDS(context.Background(), s.etcd, names)
	return err
}

func (s *DSStore) checkListener() error {
	ctx, cancel := context.WithTimeout(context.Background(), readyTO)
	defer cancel()

	err := s.listener.Ready(ctx)
	return err
}
