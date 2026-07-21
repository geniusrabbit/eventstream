package source

import (
	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream/pkg/converter"
)

// Options configures a shared source subscriber.
type Options struct {
	// Format is the payload decoder (json, bson, raw).
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

// Option mutates Options.
type Option func(opts *Options)

// WithFormat selects a named payload converter for the source.
func WithFormat(format string) Option {
	return func(opts *Options) {
		opts.Format = converter.ByName(format)
	}
}
