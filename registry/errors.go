package registry

import (
	"fmt"
)

type ErrNameConflict struct {
	conflict string
}

func (e ErrNameConflict) Error() string {
	return fmt.Sprintf("duplicate name %v found", e.conflict)
}

type ErrNotFound struct {
	notfound string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%v has not been registered", e.notfound)
}
