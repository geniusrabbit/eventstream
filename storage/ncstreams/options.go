package ncstreams

// Options contains options of the source
type Options struct {
	// Debug mode of the source
	Debug bool
}

// Option modificator callback
type Option func(opts *Options)

// WithDebug puts debug mode into the source
func WithDebug(debug bool) Option {
	return func(opts *Options) {
		opts.Debug = debug
	}
}
