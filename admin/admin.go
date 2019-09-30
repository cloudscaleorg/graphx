package admin

import (
	"github.com/cloudscaleorg/graphx"
	"github.com/cloudscaleorg/graphx/etcd"
	"github.com/cloudscaleorg/graphx/registry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// All aggregates all interfaces for administoring a graphx cluster
type All interface {
	DataSource
	Chart
}

// DataSource exports methods for administoring DataSource resources
type DataSource interface {
	CreateDataSource([]*graphx.DataSource) error
	ReadDataSource() ([]*graphx.DataSource, error)
	UpdateDataSource(ds *graphx.DataSource) error
	DeleteDataSource(ds *graphx.DataSource) error
}

// Chart exports methods for administoring Chart resources
type Chart interface {
	CreateChart([]*graphx.Chart) error
	ReadChart() ([]*graphx.Chart, error)
	ReadChartsByName(names []string) ([]*graphx.Chart, error)
	UpdateChart(ds *graphx.Chart) error
	DeleteChart(ds *graphx.Chart) error
}

// admin implements the All interface
type admin struct {
	dsmap    *etcd.DSMap
	chartmap *etcd.ChartMap
	reg      registry.Querier
	logger   zerolog.Logger
}

func NewAdmin(dsmap *etcd.DSMap, chartmap *etcd.ChartMap, reg registry.Querier) All {
	return &admin{
		dsmap:    dsmap,
		chartmap: chartmap,
		reg:      reg,
		logger:   log.With().Str("component", "admin").Logger(),
	}
}
