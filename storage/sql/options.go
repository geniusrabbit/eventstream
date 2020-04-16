package sql

import (
	"time"
)

// Option type of the stream
type Option func(stream *StreamSQL) error

// WithDebug setup the debug mode ON/OFF
func WithDebug(debug bool) Option {
	return func(stream *StreamSQL) error {
		stream.debug = debug
		return nil
	}
}

// WithBlockSize setup the size of writing block
func WithBlockSize(size int) Option {
	return func(stream *StreamSQL) error {
		stream.blockSize = size
		return nil
	}
}

// WithFlushIntervals setup interval betwin flashes
func WithFlushIntervals(interval time.Duration) Option {
	return func(stream *StreamSQL) error {
		stream.flushInterval = interval
		return nil
	}
}

// WithQueryObject setup query object
func WithQueryObject(query *Query) Option {
	return func(stream *StreamSQL) error {
		stream.query = query
		return nil
	}
}

// WithQueryRawFields setup query object by fields parameters
func WithQueryRawFields(query string, fields interface{}) Option {
	return func(stream *StreamSQL) error {
		q, err := NewQueryByRaw(query, fields)
		stream.query = q
		return err
	}
}

// WithQueryByPattern setup query object by query pattern
func WithQueryByPattern(pattern, target string, fields interface{}) Option {
	return func(stream *StreamSQL) error {
		query, err := NewQueryByPattern(pattern, target, fields)
		stream.query = query
		return err
	}
}
