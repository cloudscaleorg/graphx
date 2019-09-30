package admin

import (
	"github.com/cloudscaleorg/graphx"
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
	UpdateChart(ds *graphx.Chart) error
	DeleteChart(ds *graphx.Chart) error
}

// admin implements the All interface
type admin struct {
	dsStore    graphx.DataSourceStore
	chartStore graphx.ChartStore
	reg        registry.Querier
	logger     zerolog.Logger
}

func NewAdmin(dsStore graphx.DataSourceStore, chartStore graphx.ChartStore, reg registry.Querier) All {
	return &admin{
		dsStore:    dsStore,
		chartStore: chartStore,
		reg:        reg,
		logger:     log.With().Str("component", "admin").Logger(),
	}
}
