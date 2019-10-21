package admin

import (
	"fmt"
)

type ErrNotFound struct {
	Resource string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("resoure %v not found", e.Resource)
}

// ErrStore indicates an issue with a provided Store.
type ErrStore struct {
	error
}

// ErrMissingDataSources indicates a chart is attempted to be created with a missing datasource resource.
type ErrMissingDataSources struct {
	Missing []string
}

func (e ErrMissingDataSources) Error() string {
	return fmt.Sprintf("missing datasources: %v", e.Missing)
}

// ErrMissingBackend indicates a backend has not been implemented.
type ErrMissingBackends struct {
	Missing []string
}

func (e ErrMissingBackends) Error() string {
	return fmt.Sprintf("missing backends: %v", e.Missing)
}
