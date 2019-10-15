package storage

import (
	"encoding/json"

	"github.com/geniusrabbit/eventstream"
)

// Config of the storage
type Config = eventstream.StorageConfig

// Option of the connection
type Option func(cnf *Config)

// WithConfig custom
func WithConfig(conf *Config) Option {
	return func(cnf *Config) {
		*cnf = *conf
	}
}

// WithDebug mode
func WithDebug(debug bool) Option {
	return func(cnf *Config) {
		cnf.Debug = debug
	}
}

// WithConnect to the database
func WithConnect(driver, connect string) Option {
	return func(cnf *Config) {
		cnf.Driver = driver
		cnf.Connect = connect
	}
}

// WithBuffer size of the stream
func WithBuffer(size uint) Option {
	return func(cnf *Config) {
		cnf.Buffer = size
	}
}

// WithRawConfig storage config
func WithRawConfig(raw json.RawMessage) Option {
	return func(cnf *Config) {
		cnf.Raw = raw
	}
}

// WithObjectConfig converts Object to JSON storage config
func WithObjectConfig(obj interface{}) Option {
	return func(cnf *Config) {
		data, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		cnf.Raw = json.RawMessage(data)
	}
}
