//
// @project geniusrabbit::eventstream 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package source

import "github.com/geniusrabbit/eventstream"

var _registry = registry{
	close:      make(chan bool, 1),
	sources:    map[string]eventstream.Sourcer{},
	connectors: map[string]connector{},
}

// RegisterConnector function which creates new stream coneection by config
// and bind the connector to the `driver` name, in the global registry
func RegisterConnector(c connector, driver string) {
	_registry.RegisterConnector(c, driver)
}

// Register data source connection with `name` and options
// The source defines by connection options, in the global registry
func Register(name string, options ...Option) error {
	return _registry.Register(name, options...)
}

// Subscribe some handler interface to processing the stream with `name`,
// in the global registry
func Subscribe(name string, stream eventstream.Streamer) error {
	return _registry.Subscribe(name, stream)
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
func Listen() error {
	return _registry.Listen()
}
