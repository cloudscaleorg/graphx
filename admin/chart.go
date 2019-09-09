package admin

import "github.com/cloudscaleorg/graphx"

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

func (a *admin) CreateChart(charts []*graphx.Chart) error {
	_, missing, err := a.dsStore.GetByNames(sources(charts))
	if err != nil {
		return ErrStore{err}
	}
	if len(missing) > 0 {
		return ErrMissingDataSources{missing}
	}

	err = a.chartStore.Store(charts)
	return ErrStore{err}
}

func (a *admin) ReadChart() ([]*graphx.Chart, error) {
	sources, err := a.chartStore.Get()
	return sources, ErrStore{err}
}

func (a *admin) UpdateChart(chart *graphx.Chart) error {
	charts, _, err := a.chartStore.GetByNames([]string{chart.Name})
	if err != nil {
		return ErrStore{err}
	}
	if len(charts) <= 0 {
		return ErrNotFound
	}

	_, missing, err := a.dsStore.GetByNames(sources([]*graphx.Chart{chart}))
	if err != nil {
		return ErrStore{err}
	}
	if len(missing) > 0 {
		return ErrMissingDataSources{missing}
	}

	err = a.chartStore.Store([]*graphx.Chart{chart})
	if err != nil {
		return ErrStore{err}
	}

	return nil
}

func (a *admin) DeleteChart(ds *graphx.Chart) error {
	source, _, err := a.chartStore.GetByNames([]string{ds.Name})
	if err != nil {
		return ErrStore{err}
	}
	if len(source) <= 0 {
		return ErrNotFound
	}

	err = a.chartStore.RemoveByNames([]string{ds.Name})
	if err != nil {
		return ErrStore{err}
	}

	return nil
}
