package admin

import "github.com/cloudscaleorg/graphx"

func (a *admin) CreateDataSource(sources []*graphx.DataSource) error {
	err := a.dsStore.Store(sources)
	return ErrStore{err}
}

func (a *admin) ReadDataSource() ([]*graphx.DataSource, error) {
	sources, err := a.dsStore.Get()
	return sources, ErrStore{err}
}

func (a *admin) UpdateDataSource(ds *graphx.DataSource) error {
	source, err := a.dsStore.GetByNames([]string{ds.Name})
	if err != nil {
		return ErrStore{err}
	}
	if len(source) <= 0 {
		return ErrNotFound
	}

	// overwrite
	err = a.dsStore.Store([]*graphx.DataSource{ds})
	if err != nil {
		return ErrStore{err}
	}

	return nil
}

func (a *admin) DeleteDataSource(ds *graphx.DataSource) error {
	source, err := a.dsStore.GetByNames([]string{ds.Name})
	if err != nil {
		return ErrStore{err}
	}
	if len(source) <= 0 {
		return ErrNotFound
	}

	err = a.dsStore.RemoveByNames([]string{ds.Name})
	if err != nil {
		return ErrStore{err}
	}

	return nil
}
