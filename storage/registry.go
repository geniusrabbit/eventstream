//
// @project geniusrabbit::eventstream 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package storage

import (
	"fmt"

	"github.com/geniusrabbit/eventstream"
)

type connector func(config eventstream.ConfigItem, debug bool) (eventstream.Storager, error)

type registry struct {
	connections map[string]eventstream.Storager
	connectors  map[string]connector
}

// RegisterConnector function
func (r *registry) RegisterConnector(c connector, driver string) {
	r.connectors[driver] = c
}

// Register connection
func (r *registry) Register(name string, config eventstream.ConfigItem, debug bool) (err error) {
	var storage eventstream.Storager
	if storage, err = r.connection(config, debug); nil == err {
		r.connections[name] = storage
	}
	return
}

// Storage connection object
func (r *registry) Storage(name string) eventstream.Storager {
	conn, _ := r.connections[name]
	return conn
}

// Close listener
func (r *registry) Close() (err error) {
	for _, conn := range r.connections {
		err = conn.Close()
	}
	return
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (r *registry) connection(config eventstream.ConfigItem, debug bool) (eventstream.Storager, error) {
	var driver = config.String("driver", "")
	if conn, _ := r.connectors[driver]; conn != nil {
		return conn(config, debug)
	}
	return nil, fmt.Errorf("Undefined storage driver: [%s]", driver)
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
func Register(name string, config eventstream.ConfigItem, debug bool) error {
	return _registry.Register(name, config, debug)
}

// Storage connection object
func Storage(name string) eventstream.Storager {
	return _registry.Storage(name)
}

// Close listener
func Close() error {
	return _registry.Close()
}
