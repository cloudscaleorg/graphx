package admin

import "errors"

var (
	ErrNotFound = errors.New("resource not found")
)

// ErrStore indicates an issue with a provided Store.
type ErrStore struct {
	error
}
