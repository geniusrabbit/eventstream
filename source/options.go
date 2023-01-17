package source

import (
	"encoding/json"
	"net/url"
)

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

// WithConnect to the source
func WithConnect(driver, connect string) Option {
	return func(cnf *Config) {
		cnf.Driver = driver
		cnf.Connect = connect
	}
}

// WithConnectURL to the source
func WithConnectURL(connect string) Option {
	return func(cnf *Config) {
		url, err := url.Parse(connect)
		if err != nil {
			panic(err)
		}
		cnf.Driver = url.Scheme
		cnf.Connect = connect
	}
}

// WithFormat of the messages
func WithFormat(format string) Option {
	return func(cnf *Config) {
		cnf.Format = format
	}
}

// WithRawConfig storage config
func WithRawConfig(raw json.RawMessage) Option {
	return func(cnf *Config) {
		cnf.Raw = raw
	}
}

// WithObjectConfig converts Object to JSON storage config
func WithObjectConfig(obj any) Option {
	return func(cnf *Config) {
		data, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		cnf.Raw = json.RawMessage(data)
	}
}
