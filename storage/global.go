package storage

import "github.com/geniusrabbit/eventstream"

var _registry = registry{
	connections: map[string]eventstream.Storager{},
	connectors:  map[string]connector{},
}

// RegisterConnector function
func RegisterConnector(conn connector, driver string) {
	_registry.RegisterConnector(conn, driver)
}

// Register connection
func Register(name string, options ...Option) error {
	return _registry.Register(name, options...)
}

// Storage connection object
func Storage(name string) eventstream.Storager {
	return _registry.Storage(name)
}

// Close listener
func Close() error {
	return _registry.Close()
}
