//
// @project geniusrabbit::eventstream 2017, 2020
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020
//

package source

import (
	"context"
	"fmt"
	"sync"

	"github.com/geniusrabbit/eventstream"
)

type connector func(config *Config) (eventstream.Sourcer, error)

type registry struct {
	mx         sync.RWMutex
	close      chan bool
	connectors map[string]connector
	sources    map[string]eventstream.Sourcer
}

// RegisterConnector function which creates new stream coneection by config
// and bind the connector to the `driver` name
func (r *registry) RegisterConnector(c connector, driver string) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.connectors[driver] = c
}

// Register data source connection with `name` and options
// The source defines by connection options
func (r *registry) Register(name string, options ...Option) (err error) {
	var (
		source eventstream.Sourcer
		config Config
	)
	for _, opt := range options {
		opt(&config)
	}
	if source, err = r.connection(&config); err == nil {
		r.mx.Lock()
		defer r.mx.Unlock()
		r.sources[name] = source
	}
	return
}

// Subscribe some handler interface to processing the stream with `name`
func (r *registry) Subscribe(ctx context.Context, name string, stream eventstream.Streamer) error {
	r.mx.RLock()
	defer r.mx.RUnlock()
	if src, _ := r.sources[name]; src != nil {
		return src.Subscribe(ctx, stream)
	}
	return nil
}

// Source returns the source object registered with the `name`
func (r *registry) Source(name string) eventstream.Sourcer {
	r.mx.RLock()
	defer r.mx.RUnlock()
	src, _ := r.sources[name]
	return src
}

// Listen method launch into the background all sources where the supervised
// daemon mode is required
func (r *registry) Listen(ctx context.Context) (err error) {
	r.mx.RLock()
	for _, source := range r.sources {
		if err = source.Start(ctx); err != nil {
			r.Close()
			r.mx.RUnlock()
			return err
		}
	}
	r.mx.RUnlock()
	<-r.close
	return
}

// Close all listeners and source connections
func (r *registry) Close() (err error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	for _, source := range r.sources {
		err = source.Close()
	}
	r.close <- true
	return
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (r *registry) connection(config *Config) (eventstream.Sourcer, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	if conn, _ := r.connectors[config.Driver]; conn != nil {
		return conn(config)
	}
	return nil, fmt.Errorf("[source] undefined source driver: [%s]", config.Driver)
}
