package ncstreams

import (
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream/converter"
)

// Options contains options of the source
type Options struct {
	// Debug mode of the source
	Debug bool

	// Format converter
	Format converter.Converter
}

func (opts *Options) getFormat() converter.Converter {
	if opts.Format != nil {
		return opts.Format
	}
	return converter.ByName(`raw`)
}

func (opts *Options) getLogger() *zap.Logger {
	return zap.L()
}

// Option modificator callback
type Option func(opts *Options)

// WithDebug puts debug mode into the source
func WithDebug(debug bool) Option {
	return func(opts *Options) {
		opts.Debug = debug
	}
}

// WithFormat puts format name into the source
func WithFormat(format string) Option {
	return func(opts *Options) {
		opts.Format = converter.ByName(format)
	}
}
