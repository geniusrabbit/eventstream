package sql

import (
	"time"
)

// Option type of the stream
type Option func(stream *StreamSQL)

// WithDebug setup the debug mode ON/OFF
func WithDebug(debug bool) Option {
	return func(stream *StreamSQL) {
		stream.debug = debug
	}
}

// WithBlockSize setup the size of writing block
func WithBlockSize(size int) Option {
	return func(stream *StreamSQL) {
		stream.blockSize = size
	}
}

// WithFlushIntervals setup interval betwin flashes
func WithFlushIntervals(interval time.Duration) Option {
	return func(stream *StreamSQL) {
		stream.flushInterval = interval
	}
}
