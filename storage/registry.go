//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package storage

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream"
)

// Error list...
var (
	ErrUndefinedDriver = errors.New(`[storage::registry] undefined driver`)
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
func (r *registry) Register(name string, options ...Option) (err error) {
	var (
		storage eventstream.Storager
		config  Config
	)
	for _, opt := range options {
		opt(&config)
	}
	if storage, err = r.connection(&config); err == nil {
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
	return nil, errors.Wrap(ErrUndefinedDriver, config.Driver)
}
