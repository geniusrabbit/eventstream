//
// @project geniusrabbit::eventstream 2017, 2020 - 2023
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2020 - 2023
//

package source

import (
	"context"

	"github.com/geniusrabbit/eventstream"
)

var _registry = registry{
	close:      make(chan bool, 1),
	sources:    map[string]eventstream.Sourcer{},
	connectors: map[string]connector{},
}

// RegisterConnector function which creates new stream coneection by config
// and bind the connector to the `driver` name, in the global registry
func RegisterConnector(driver string, c connector) {
	_registry.RegisterConnector(driver, c)
}

// Register data source connection with `name` and options
// The source defines by connection options, in the global registry
func Register(ctx context.Context, name string, options ...Option) error {
	return _registry.Register(ctx, name, options...)
}

// Subscribe some handler interface to processing the stream with `name`,
// in the global registry
func Subscribe(ctx context.Context, name string, streams ...eventstream.Streamer) error {
	return _registry.Subscribe(ctx, name, streams...)
}

// Source returns the source object registered with the `name`, in the global registry
func Source(name string) eventstream.Sourcer {
	return _registry.Source(name)
}

// Close all listeners and source connections, in the global registry
func Close() error {
	return _registry.Close()
}

// Listen method launch into the background all sources where the supervised
// daemon mode is required, in the global registry
func Listen(ctx context.Context) error {
	return _registry.Listen(ctx)
}
