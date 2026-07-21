package pipeline

import (
	"context"

	"go.uber.org/zap"

	"github.com/geniusrabbit/eventstream"
)

// Option configures an Engine during New. Receives the same ctx passed to New.
type Option interface {
	apply(ctx context.Context, e *Engine) error
	phase() int
}

const (
	phaseConfig  = 0
	phaseSource  = 1
	phaseStorage = 2
	phaseStream  = 3
)

type optionFunc struct {
	p  int
	fn func(ctx context.Context, e *Engine) error
}

func (o optionFunc) apply(ctx context.Context, e *Engine) error { return o.fn(ctx, e) }
func (o optionFunc) phase() int                                 { return o.p }

// WithDebug enables debug flags on the engine.
func WithDebug(debug bool) Option {
	return optionFunc{p: phaseConfig, fn: func(_ context.Context, e *Engine) error {
		e.Debug = debug
		return nil
	}}
}

// WithLogger sets the engine logger.
func WithLogger(logger *zap.Logger) Option {
	return optionFunc{p: phaseConfig, fn: func(_ context.Context, e *Engine) error {
		e.Logger = logger
		return nil
	}}
}

// Source registers an already opened source.
func Source(name string, src eventstream.Sourcer) Option {
	return optionFunc{p: phaseSource, fn: func(ctx context.Context, e *Engine) error {
		return e.RegisterSource(ctx, name, src)
	}}
}

// SourceFunc opens and registers a source using New's context.
func SourceFunc(name string, fn func(ctx context.Context) (eventstream.Sourcer, error)) Option {
	return optionFunc{p: phaseSource, fn: func(ctx context.Context, e *Engine) error {
		src, err := fn(ctx)
		if err != nil {
			return err
		}
		return e.RegisterSource(ctx, name, src)
	}}
}

// Storage registers an already opened storage.
func Storage(name string, store eventstream.Storager) Option {
	return optionFunc{p: phaseStorage, fn: func(ctx context.Context, e *Engine) error {
		return e.RegisterStorage(ctx, name, store)
	}}
}

// StorageFunc opens and registers a storage using New's context.
func StorageFunc(name string, fn func(ctx context.Context) (eventstream.Storager, error)) Option {
	return optionFunc{p: phaseStorage, fn: func(ctx context.Context, e *Engine) error {
		store, err := fn(ctx)
		if err != nil {
			return err
		}
		return e.RegisterStorage(ctx, name, store)
	}}
}

// Stream registers a stream binding (applied after sources/storages in New).
func Stream(name string, opts ...StreamOption) Option {
	return optionFunc{p: phaseStream, fn: func(ctx context.Context, e *Engine) error {
		return e.RegisterStream(ctx, name, opts...)
	}}
}
