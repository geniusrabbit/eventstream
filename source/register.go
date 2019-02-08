//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package source

import (
	"fmt"
	"sync"

	"github.com/geniusrabbit/eventstream"
)

type connector func(config eventstream.ConfigItem, debug bool) (eventstream.Sourcer, error)

type registry struct {
	mx         sync.RWMutex
	close      chan bool
	connectors map[string]connector
	sources    map[string]eventstream.Sourcer
}

// RegisterConnector stream subscriber
func (r *registry) RegisterConnector(c connector, driver string) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.connectors[driver] = c
}

// Register stream subscriber
func (r *registry) Register(name string, config eventstream.ConfigItem, debug bool) (err error) {
	var source eventstream.Sourcer
	if source, err = r.connection(config, debug); err == nil {
		r.mx.Lock()
		defer r.mx.Unlock()
		r.sources[name] = source
	}
	return
}

// Subscribe handler
func (r *registry) Subscribe(name string, stream eventstream.Streamer) error {
	r.mx.RLock()
	defer r.mx.RUnlock()
	if src, _ := r.sources[name]; src != nil {
		return src.Subscribe(stream)
	}
	return nil
}

// Source object by name
func (r *registry) Source(name string) eventstream.Sourcer {
	r.mx.RLock()
	defer r.mx.RUnlock()
	src, _ := r.sources[name]
	return src
}

// Listen sources
func (r *registry) Listen() (err error) {
	r.mx.RLock()
	for _, source := range r.sources {
		if err = source.Start(); err != nil {
			r.Close()
			r.mx.RUnlock()
			return err
		}
	}
	r.mx.RUnlock()
	<-r.close
	return
}

// Close listener
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

func (r *registry) connection(config eventstream.ConfigItem, debug bool) (eventstream.Sourcer, error) {
	var driver = config.String("driver", "")

	r.mx.RLock()
	defer r.mx.RUnlock()
	if conn, _ := r.connectors[driver]; conn != nil {
		return conn(config, debug)
	}
	return nil, fmt.Errorf("Undefined source driver: [%s]", driver)
}

///////////////////////////////////////////////////////////////////////////////
/// Global
///////////////////////////////////////////////////////////////////////////////

var _registry = registry{
	close:      make(chan bool, 1),
	sources:    map[string]eventstream.Sourcer{},
	connectors: map[string]connector{},
}

// RegisterConnector stream subscriber
func RegisterConnector(c connector, driver string) {
	_registry.RegisterConnector(c, driver)
}

// Register connection
func Register(name string, config eventstream.ConfigItem, debug bool) error {
	return _registry.Register(name, config, debug)
}

// Subscribe handler
func Subscribe(name string, stream eventstream.Streamer) error {
	return _registry.Subscribe(name, stream)
}

// Source object by name
func Source(name string) eventstream.Sourcer {
	return _registry.Source(name)
}

// Close listener
func Close() error {
	return _registry.Close()
}

// Listen sources
func Listen() error {
	return _registry.Listen()
}
