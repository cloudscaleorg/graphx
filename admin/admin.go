package admin

import "github.com/cloudscaleorg/graphx"

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
}

func NewAdmin(dsStore graphx.DataSourceStore, chartStore graphx.ChartStore) All {
	return &admin{
		dsStore:    dsStore,
		chartStore: chartStore,
	}
}
