//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package storage

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream"
)

// Error list...
var (
	ErrUndefinedDriver = errors.New(`[storage::registry] undefined driver`)
)

type connector func(ctx context.Context, config *Config) (eventstream.Storager, error)

type registry struct {
	mx          sync.RWMutex
	connections map[string]eventstream.Storager
	connectors  map[string]connector
}

// RegisterConnector function
func (r *registry) RegisterConnector(driver string, c connector) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.connectors[driver] = c
}

// Register connection
func (r *registry) Register(ctx context.Context, name string, options ...Option) (err error) {
	var (
		storage eventstream.Storager
		config  Config
	)
	for _, opt := range options {
		opt(&config)
	}
	if storage, err = r.connection(ctx, &config); err == nil {
		r.mx.Lock()
		defer r.mx.Unlock()
		r.connections[name] = storage
	}
	return err
}

// Storage connection object
func (r *registry) Storage(name string) eventstream.Storager {
	r.mx.RLock()
	defer r.mx.RUnlock()
	return r.connections[name]
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
	return err
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (r *registry) connection(ctx context.Context, config *Config) (eventstream.Storager, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	if conn := r.connectors[config.Driver]; conn != nil {
		return conn(ctx, config)
	}
	return nil, errors.Wrap(ErrUndefinedDriver, config.Connect)
}
