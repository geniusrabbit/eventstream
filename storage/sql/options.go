package sql

import (
	"time"

	"go.uber.org/zap"
)

const (
	defaultBlockSize     = 1000
	defaultFlushInterval = time.Second * 1
)

// Options of the SQL stream
type Options struct {
	// Debug mode of the stream
	Debug bool

	// BlockSize of buffer to flushing data
	BlockSize int

	// FlushInterval between flushes
	FlushInterval time.Duration

	// QueryBuilder object of data prepare
	QueryBuilder *Query

	// Logger object of module debuging and logs
	Logger *zap.Logger
}

func (opts *Options) getBlockSize() int {
	if opts.BlockSize < 1 {
		return defaultBlockSize
	}
	return opts.BlockSize
}

func (opts *Options) getFlushInterval() time.Duration {
	if opts.FlushInterval < 1 {
		return defaultFlushInterval
	}
	return opts.FlushInterval
}

func (opts *Options) getLogger() *zap.Logger {
	if opts.Logger == nil {
		return zap.L().With(zap.String(`module`, `sql-stream`))
	}
	return opts.Logger
}

// Option type of the stream
type Option func(opts *Options) error

// WithDebug setup the debug mode ON/OFF
func WithDebug(debug bool) Option {
	return func(opts *Options) error {
		opts.Debug = debug
		return nil
	}
}

// WithLogger setup the logging object
func WithLogger(logger *zap.Logger) Option {
	return func(opts *Options) error {
		opts.Logger = logger
		return nil
	}
}

// WithBlockSize setup the size of writing block
func WithBlockSize(size int) Option {
	return func(opts *Options) error {
		opts.BlockSize = size
		return nil
	}
}

// WithFlushIntervals setup interval betwin flashes
func WithFlushIntervals(interval time.Duration) Option {
	return func(opts *Options) error {
		opts.FlushInterval = interval
		return nil
	}
}

// WithQueryObject setup query object
func WithQueryObject(query *Query) Option {
	return func(opts *Options) error {
		opts.QueryBuilder = query
		return nil
	}
}

// WithQueryRawFields setup query object by fields parameters
func WithQueryRawFields(query string, fields interface{}) Option {
	return func(opts *Options) error {
		queryBuilder, err := NewQueryByRaw(query, fields)
		opts.QueryBuilder = queryBuilder
		return err
	}
}

// WithQueryByPattern setup query object by query pattern
func WithQueryByPattern(pattern, target string, fields interface{}) Option {
	return func(opts *Options) error {
		queryBuilder, err := NewQueryByPattern(pattern, target, fields)
		opts.QueryBuilder = queryBuilder
		return err
	}
}
