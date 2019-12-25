package admin

import (
	"github.com/cloudscaleorg/graphx"
)

// extracts a list of datasource names from a list of charts
func datasources(charts []*graphx.Chart) []string {
	out := []string{}
	seen := map[string]struct{}{}
	for _, chart := range charts {
		for datasource, _ := range chart.DataSources {
			if _, ok := seen[datasource]; !ok {
				out = append(out, datasource)
			}
		}
	}
	return out
}

// extracts a list of chart names from a list of charts
func names(charts []*graphx.Chart) []string {
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
	_, missing := a.dsmap.Get(datasources(charts))
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
	_, missing := a.dsmap.Get(datasources([]*graphx.Chart{chart}))
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
