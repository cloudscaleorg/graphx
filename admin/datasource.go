package admin

import (
	"github.com/cloudscaleorg/graphx"
)

func (a *admin) CreateDataSource(sources []*graphx.DataSource) error {
	missing := []string{}
	for _, ds := range sources {
		if ok := a.beReg.Exists(ds.Backend); !ok {
			missing = append(missing, ds.Backend)
		}
	}
	if len(missing) > 0 {
		return ErrMissingBackends{missing}
	}
	a.dsmap.Store(sources)
	return nil
}

func (a *admin) ReadDataSource() ([]*graphx.DataSource, error) {
	sources, _ := a.dsmap.Get(nil)
	return sources, nil
}

func (a *admin) ReadDataSourcesByName(names []string) ([]*graphx.DataSource, error) {
	datasources, _ := a.dsmap.Get(names)
	return datasources, nil
}

func (a *admin) UpdateDataSource(ds *graphx.DataSource) error {
	source, _ := a.dsmap.Get([]string{ds.Name})
	if len(source) <= 0 {
		return ErrNotFound{ds.Name}
	}
	if ok := a.beReg.Exists(ds.Backend); !ok {
		return ErrMissingBackends{[]string{ds.Backend}}
	}
	// overwrite
	a.dsmap.Store([]*graphx.DataSource{ds})
	return nil
}

func (a *admin) DeleteDataSource(ds *graphx.DataSource) error {
	source, _ := a.dsmap.Get([]string{ds.Name})
	if len(source) <= 0 {
		return ErrNotFound{ds.Name}
	}
	a.dsmap.Remove([]string{ds.Name})
	return nil
}
