package registry

import "github.com/cloudscaleorg/graphx"

type QuerierOpts struct {
	Metrics    []graphx.ChartMetric
	ConnString string
	MChan      chan *graphx.Metric
	EChan      chan error
}

type QuerierFactory func(opts QuerierOpts) graphx.Querier

type Registry interface {
	Registerer
	Getter
}

type Registerer interface {
	RegisterQuerier(name string, factory QuerierFactory) error
}
type Getter interface {
	GetQuerier(name string) (graphx.Querier, error)
}

// registry implements
type registry struct {
	qMap map[string]QuerierFactory
}

func NewRegistry() Registerer {
	return &registry{
		qMap: map[string]QuerierFactory{},
	}
}

func (r *registry) RegisterQuerier(name string, factory QuerierFactory) error {

}
