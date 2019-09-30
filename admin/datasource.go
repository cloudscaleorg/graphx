package admin

import (
	"github.com/cloudscaleorg/graphx"
)

func (a *admin) CreateDataSource(sources []*graphx.DataSource) error {
	err := a.dsStore.Store(sources)
	if err != nil {
		return ErrStore{err}
	}
	return nil
}

func (a *admin) ReadDataSource() ([]*graphx.DataSource, error) {
	sources, err := a.dsStore.Get()
	if err != nil {
		return nil, ErrStore{err}
	}
	return sources, nil
}

func (a *admin) UpdateDataSource(ds *graphx.DataSource) error {
	source, _, err := a.dsStore.GetByNames([]string{ds.Name})
	if err != nil {
		return ErrStore{err}
	}
	if len(source) <= 0 {
		return ErrNotFound{ds.Name}
	}

	// overwrite
	err = a.dsStore.Store([]*graphx.DataSource{ds})
	if err != nil {
		return ErrStore{err}
	}

	return nil
}

func (a *admin) DeleteDataSource(ds *graphx.DataSource) error {
	source, _, err := a.dsStore.GetByNames([]string{ds.Name})
	if err != nil {
		return ErrStore{err}
	}
	if len(source) <= 0 {
		return ErrNotFound{ds.Name}
	}

	err = a.dsStore.RemoveByNames([]string{ds.Name})
	if err != nil {
		return ErrStore{err}
	}

	return nil
}
