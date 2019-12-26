package admin

import (
	"github.com/cloudscaleorg/graphx"
)

// extracts a list of datasource names from a list of Charts
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

// extracts a list of Chart names from a list of Charts
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

// CreateChart persists a slice of Charts.
//
// If a provided Chart contains any references to a DataSource
// that has not been created prior an ErrMissingDataSources error is returned
func (a *Admin) CreateChart(charts []*graphx.Chart) error {
	_, missing := a.dsmap.Get(datasources(charts))
	if len(missing) > 0 {
		return ErrMissingDataSources{missing}
	}
	a.chartmap.Store(charts)
	return nil
}

// ReadChart lists all created Charts
func (a *Admin) ReadChart() ([]*graphx.Chart, error) {
	sources, _ := a.chartmap.Get(nil)
	return sources, nil
}

// ReadChartsByName lists Charts by their unique names
func (a *Admin) ReadChartsByName(names []string) ([]*graphx.Chart, error) {
	charts, _ := a.chartmap.Get(names)
	return charts, nil
}

// UpdateChart first confirms the provided Chart exists
// and if so overwrites the original with the provided.
//
// If the provided Chart does not exist an ErrNotFound
// error will be returned
//
// If a provided Chart contains any references to a DataSource
// that has not been created prior an ErrMissingDataSources error is returned
func (a *Admin) UpdateChart(chart *graphx.Chart) error {
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

// DeleteChart removes the existence of the provided Chart
//
// If the provided Chart does not exist an ErrNotFound
// error will be returned
func (a *Admin) DeleteChart(ds *graphx.Chart) error {
	source, _ := a.chartmap.Get([]string{ds.Name})
	if len(source) <= 0 {
		return ErrNotFound{ds.Name}
	}
	a.chartmap.Remove([]string{ds.Name})
	return nil
}
