package admin

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("resource not found")
)

// ErrStore indicates an issue with a provided Store.
type ErrStore struct {
	error
}

type ErrMissingDataSources struct {
	missing []string
}

func (e ErrMissingDataSources) Error() string {
	return fmt.Sprintf("missing datasources: %v", e.missing)
}
