package admin

import (
	"github.com/cloudscaleorg/graphx"
)

// CreateDataSource persists a slice of DataSources.
//
// If a provided DataSource contains any references to a Backend
// that has not been implemented an ErrMissingBackends error is returned
func (a *Admin) CreateDataSource(sources []*graphx.DataSource) error {
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

// ReadDataSource lists all created DataSources
func (a *Admin) ReadDataSource() ([]*graphx.DataSource, error) {
	sources, _ := a.dsmap.Get(nil)
	return sources, nil
}

// ReadDataSourcesByName lists DataSources by their unique names
func (a *Admin) ReadDataSourcesByName(names []string) ([]*graphx.DataSource, error) {
	datasources, _ := a.dsmap.Get(names)
	return datasources, nil
}

// UpdateDataSource first confirms the provided DataSource exists
// and if so overwrites the original with the provided.
//
// If the provided DataSource does not exist an ErrNotFound
// error will be returned
//
// If a provided DataSource contains any references to a Backend
// not implemented an ErrMissingBackends error is returned
func (a *Admin) UpdateDataSource(ds *graphx.DataSource) error {
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

// DeleteDataSource removes the existence of the provided DataSource
//
//
// If a provided DataSource contains any references to a Backend
// not implemented an ErrMissingBackends error is returned
func (a *Admin) DeleteDataSource(ds *graphx.DataSource) error {
	source, _ := a.dsmap.Get([]string{ds.Name})
	if len(source) <= 0 {
		return ErrNotFound{ds.Name}
	}
	a.dsmap.Remove([]string{ds.Name})
	return nil
}
