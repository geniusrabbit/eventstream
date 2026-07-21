package storage

import "github.com/pkg/errors"

var errInvalidMetricsItemConfig = errors.New("[storage] invalid stream targets config")

// Target describes one published message shape.
type Target struct {
	Fields map[string]any
	Where  string
}

// StreamOption configures a publish sink stream.
type StreamOption func(*streamOptions)

type streamOptions struct {
	id      string
	targets []Target
}

// WithID sets the stream identifier.
func WithID(id string) StreamOption {
	return func(o *streamOptions) { o.id = id }
}

// WithTargets sets publish targets.
func WithTargets(targets ...Target) StreamOption {
	return func(o *streamOptions) { o.targets = append(o.targets, targets...) }
}
