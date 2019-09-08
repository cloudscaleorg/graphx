package etcd

import "fmt"

type ErrNotFound struct {
	missing []string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("following charts not found: %v", e.missing)
}
