package admin

import (
	"github.com/cloudscaleorg/graphx"
)

// retrieves a list of datasource names from a list of charts
var sources = func(charts []*graphx.Chart) []string {
	seen := map[string]struct{}{}
	for _, chart := range charts {
		for _, metric := range chart.Metrics {
			seen[metric.DataSource] = struct{}{}
		}
	}
	out := []string{}
	for k, _ := range seen {
		out = append(out, k)
	}
	return out
}

// names retrieves a list of chart names from a list of charts
var names = func(charts []*graphx.Chart) []string {
	seen := map[string]struct{}{}
	names := []string{}
	for _, chart := range charts {
		if _, ok := seen[chart.Name]; !ok {
			names = append(names, chart.Name)
		}
	}

	return names
}

func (a *admin) CreateChart(charts []*graphx.Chart) error {
	_, missing := a.dsmap.Get(sources(charts))
	if len(missing) > 0 {
		return ErrMissingDataSources{missing}
	}

	a.chartmap.Store(charts)
	return nil
}

func (a *admin) ReadChart() ([]*graphx.Chart, error) {
	sources, _ := a.chartmap.Get(nil)
	return sources, nil
}

func (a *admin) ReadChartsByName(names []string) ([]*graphx.Chart, error) {
	charts, _ := a.chartmap.Get(names)
	return charts, nil
}

func (a *admin) UpdateChart(chart *graphx.Chart) error {
	charts, _ := a.chartmap.Get([]string{chart.Name})
	if len(charts) <= 0 {
		return ErrNotFound{chart.Name}
	}

	_, missing := a.dsmap.Get(sources([]*graphx.Chart{chart}))
	if len(missing) > 0 {
		return ErrMissingDataSources{missing}
	}

	a.chartmap.Store([]*graphx.Chart{chart})
	return nil
}

func (a *admin) DeleteChart(ds *graphx.Chart) error {
	source, _ := a.chartmap.Get([]string{ds.Name})
	if len(source) <= 0 {
		return ErrNotFound{ds.Name}
	}

	a.chartmap.Remove([]string{ds.Name})
	return nil
}

func (a *admin) ChartMetricsByDataSource(names []string) map[*graphx.DataSource][]*graph.ChartMetrics {
	charts, missing := a.chartmap.Get(names)
	for _, chart := range charts {
		for chart.Metrics {
			
		}
	}


	datasources := map[string][]graphx.ChartMetrics{}
	out := map[*graphx.DataSource][]*graphx.ChartMetics{}

	for _, chart := range chart {
		if _, !ok := datasources[chart.Name]
	}
}
