package admin

import (
	"github.com/cloudscaleorg/graphx"
)

func (a *admin) CreateDataSource(sources []*graphx.DataSource) error {
	missing := []string{}
	for _, ds := range sources {
		if ok := a.reg.Check(ds.Type); !ok {
			missing = append(missing, ds.Type)
		}
	}

	if len(missing) > 0 {
		return ErrMissingQueriers{missing}
	}

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
	if ok := a.reg.Check(ds.Type); !ok {
		return ErrMissingQueriers{[]string{ds.Type}}
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
