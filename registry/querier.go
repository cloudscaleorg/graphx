package registry

import "github.com/cloudscaleorg/graphx"

type QuerierOpts struct {
	Metrics    []graphx.ChartMetric
	ConnString string
	MChan      chan *graphx.Metric
	EChan      chan error
}

type QuerierFactory func(opts QuerierOpts) graphx.Querier

type Querier interface {
	Registerer
	Getter
	Checker
}

type Registerer interface {
	Register(name string, factory QuerierFactory) error
}
type Getter interface {
	Get(name string, opts QuerierOpts) (graphx.Querier, error)
}
type Checker interface {
	Check(name string) bool
}

// registry implements the registry.Querier interface
type registry struct {
	qMap map[string]QuerierFactory
}

func NewRegistry() Registerer {
	return &registry{
		qMap: map[string]QuerierFactory{},
	}
}

func (r *registry) Register(name string, factory QuerierFactory) error {
	if _, ok := r.qMap[name]; ok {
		return ErrNameConflict{name}
	}

	r.qMap[name] = factory
	return nil
}

func (r *registry) Get(name string, opts QuerierOpts) (graphx.Querier, error) {
	factory, ok := r.qMap[name]
	if !ok {
		return nil, ErrNotFound{name}
	}

	return factory(opts), nil
}

func (r *registry) Check(name string) bool {
	_, ok := r.qMap[name]
	return ok
}
