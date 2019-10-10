package registry

import (
	"sync"

	"github.com/cloudscaleorg/graphx"
)

// BackendFactory are a function type used in registration.
//
// the factory should return a *graphx.Backend able to connect to the datasource provided
type BackendFactory func(ds *graphx.DataSource) (*graphx.Backend, error)

// Backend is a registry for graphx.Backend implementations.
//
// acts as a factory for creating configured Backends and checking if a particular backend exists
type Backend interface {
	// Register a backend by name. its an error to register two or more backend of of the same name
	Register(name string, f BackendFactory) error
	// Get must take a datasource and return a *graphx.Backend or an error if the backend does not exist
	Get(ds *graphx.DataSource) (*graphx.Backend, error)
	// Exists must tell the caller if the backend exists in the registry given a name
	Exists(name string) bool
}

// bregistry implements the Backend interface
type bregistry struct {
	mu  sync.RWMutex
	reg map[string]BackendFactory
}

func (r *bregistry) Register(name string, f BackendFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.reg[name]; ok {
		return &ErrDuplicateBackend{name}
	}
	r.reg[name] = f
	return nil
}

func (r *bregistry) Get(ds *graphx.DataSource) (*graphx.Backend, error) {
	var (
		factory BackendFactory
		ok      bool
	)
	r.mu.RLock()
	if factory, ok = r.reg[ds.Backend]; !ok {
		r.mu.RUnlock()
		return nil, &ErrBackendNotExist{ds.Backend}
	}
	r.mu.RUnlock()
	return factory(ds)
}

func (r *bregistry) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RLock()
	_, ok := r.reg[name]
	return ok
}
