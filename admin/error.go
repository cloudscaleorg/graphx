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

// ErrMissingDataSources indicates a chart is attempted to be created with a missing datasource resource
type ErrMissingDataSources struct {
	Missing []string
}

func (e ErrMissingDataSources) Error() string {
	return fmt.Sprintf("missing datasources: %v", e.Missing)
}

// ErrMissingQueriers indicates a querier has not been implemented
type ErrMissingQueriers struct {
	Missing []string
}

func (e ErrMissingQueriers) Error() string {
	return fmt.Sprintf("missing queriers: %v", e.Missing)
}
