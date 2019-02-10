//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package storage

import (
	"fmt"
	"sync"

	"github.com/geniusrabbit/eventstream"
)

type connector func(config *Config) (eventstream.Storager, error)

type registry struct {
	mx          sync.RWMutex
	connections map[string]eventstream.Storager
	connectors  map[string]connector
}

// RegisterConnector function
func (r *registry) RegisterConnector(c connector, driver string) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.connectors[driver] = c
}

// Register connection
func (r *registry) Register(name string, config *Config) (err error) {
	var storage eventstream.Storager
	if storage, err = r.connection(config); err == nil {
		r.mx.Lock()
		defer r.mx.Unlock()
		r.connections[name] = storage
	}
	return
}

// Storage connection object
func (r *registry) Storage(name string) eventstream.Storager {
	r.mx.RLock()
	defer r.mx.RUnlock()
	conn, _ := r.connections[name]
	return conn
}

// Close listener
func (r *registry) Close() (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, conn := range r.connections {
		err = conn.Close()
	}

	// Reset connections
	r.connections = map[string]eventstream.Storager{}
	return
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (r *registry) connection(config *Config) (eventstream.Storager, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	if conn, _ := r.connectors[config.Driver]; conn != nil {
		return conn(config)
	}
	return nil, fmt.Errorf("[storage::registry] undefined driver: [%s]", config.Driver)
}

///////////////////////////////////////////////////////////////////////////////
/// Global
///////////////////////////////////////////////////////////////////////////////

var _registry = registry{
	connections: map[string]eventstream.Storager{},
	connectors:  map[string]connector{},
}

// RegisterConnector function
func RegisterConnector(conn connector, driver string) {
	_registry.RegisterConnector(conn, driver)
}

// Register connection
func Register(name string, config *Config) error {
	return _registry.Register(name, config)
}

// Storage connection object
func Storage(name string) eventstream.Storager {
	return _registry.Storage(name)
}

// Close listener
func Close() error {
	return _registry.Close()
}
