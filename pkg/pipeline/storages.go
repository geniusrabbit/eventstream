package pipeline

import (
	"sync"

	"github.com/geniusrabbit/eventstream"
)

// storageRegistry holds named storages for one engine instance.
type storageRegistry struct {
	mx          sync.RWMutex
	connections map[string]eventstream.Storager
}

func newStorageRegistry() *storageRegistry {
	return &storageRegistry{
		connections: map[string]eventstream.Storager{},
	}
}

func (r *storageRegistry) set(name string, store eventstream.Storager) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.connections[name] = store
}

func (r *storageRegistry) get(name string) eventstream.Storager {
	r.mx.RLock()
	defer r.mx.RUnlock()
	return r.connections[name]
}

func (r *storageRegistry) closeAll() (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	for _, conn := range r.connections {
		err = conn.Close()
	}
	r.connections = map[string]eventstream.Storager{}
	return err
}
