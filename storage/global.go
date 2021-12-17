package storage

import (
	"context"

	"github.com/geniusrabbit/eventstream"
)

var _registry = registry{
	connections: map[string]eventstream.Storager{},
	connectors:  map[string]connector{},
}

// RegisterConnector function
func RegisterConnector(driver string, conn connector) {
	_registry.RegisterConnector(driver, conn)
}

// Register connection
func Register(ctx context.Context, name string, options ...Option) error {
	return _registry.Register(ctx, name, options...)
}

// Storage connection object
func Storage(name string) eventstream.Storager {
	return _registry.Storage(name)
}

// Close listener
func Close() error {
	return _registry.Close()
}
