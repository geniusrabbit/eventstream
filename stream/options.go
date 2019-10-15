package stream

import (
	"github.com/geniusrabbit/eventstream"
)

// Config of the stream
type Config = eventstream.StreamConfig

// Option of the stream
type Option func(cnf *Config)

// WithConfig custom
func WithConfig(conf *Config) Option {
	return func(cnf *Config) {
		*cnf = *conf
	}
}

// WithName of the stream
func WithName(name string) Option {
	return func(cnf *Config) {
		cnf.Name = name
	}
}

// WithDebug mode
func WithDebug(debug bool) Option {
	return func(cnf *Config) {
		cnf.Debug = debug
	}
}

// WithWhere condition
func WithWhere(where string) Option {
	return func(cnf *Config) {
		cnf.Where = where
	}
}
