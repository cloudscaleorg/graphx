package registry

import "fmt"

type ErrDuplicateBackend struct {
	backend string
}

func (e *ErrDuplicateBackend) Error() string {
	return fmt.Sprintf("attempted to register existing backend %v", e.backend)
}

type ErrBackendNotExist struct {
	backend string
}

func (e *ErrBackendNotExist) Error() string {
	return fmt.Sprintf("backend %v does not exist in registry", e.backend)
}
