package redis

import "time"

// StreamOption configures a redis sink stream.
type StreamOption func(*streamOptions)

type streamOptions struct {
	id          string
	key         string
	expiration  time.Duration
	incrementor bool
}

// WithID sets the stream identifier.
func WithID(id string) StreamOption {
	return func(o *streamOptions) { o.id = id }
}

// WithKey sets the redis key template.
func WithKey(key string) StreamOption {
	return func(o *streamOptions) { o.key = key }
}

// WithExpiration sets key TTL.
func WithExpiration(d time.Duration) StreamOption {
	return func(o *streamOptions) { o.expiration = d }
}

// WithIncrementor enables incrementor mode.
func WithIncrementor(v bool) StreamOption {
	return func(o *streamOptions) { o.incrementor = v }
}
