package admin

import "github.com/cloudscaleorg/graphx"

func (a *admin) CreateChart(sources []*graphx.Chart) error {
	err := a.chartStore.Store(sources)
	return ErrStore{err}
}

func (a *admin) ReadChart() ([]*graphx.Chart, error) {
	sources, err := a.chartStore.Get()
	return sources, ErrStore{err}
}

func (a *admin) UpdateChart(ds *graphx.Chart) error {
	source, err := a.chartStore.GetByNames([]string{ds.Name})
	if err != nil {
		return ErrStore{err}
	}
	if len(source) <= 0 {
		return ErrNotFound
	}

	// overwrite
	err = a.chartStore.Store([]*graphx.Chart{ds})
	if err != nil {
		return ErrStore{err}
	}

	return nil
}

func (a *admin) DeleteChart(ds *graphx.Chart) error {
	source, err := a.chartStore.GetByNames([]string{ds.Name})
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
