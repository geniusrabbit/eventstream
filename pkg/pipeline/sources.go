package pipeline

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/stream"
)

var errUndefinedSource = errors.New("undefined source")

// sourceRegistry holds named sources for one engine instance.
type sourceRegistry struct {
	mx      sync.RWMutex
	close   chan bool
	sources map[string]eventstream.Sourcer
}

func newSourceRegistry() *sourceRegistry {
	return &sourceRegistry{
		close:   make(chan bool, 1),
		sources: map[string]eventstream.Sourcer{},
	}
}

func (r *sourceRegistry) set(name string, src eventstream.Sourcer) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.sources[name] = src
}

func (r *sourceRegistry) subscribe(ctx context.Context, name string, streams ...stream.Streamer) error {
	r.mx.RLock()
	defer r.mx.RUnlock()
	src := r.sources[name]
	if src == nil {
		return errors.Wrap(errUndefinedSource, name)
	}
	return src.Subscribe(ctx, streams...)
}

func (r *sourceRegistry) listen(ctx context.Context) (err error) {
	r.mx.RLock()
	for _, source := range r.sources {
		if err = source.Start(ctx); err != nil {
			r.mx.RUnlock()
			_ = r.closeAll()
			return err
		}
	}
	r.mx.RUnlock()

	select {
	case <-ctx.Done():
		return nil
	case <-r.close:
		return nil
	}
}

func (r *sourceRegistry) closeAll() (err error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	for _, source := range r.sources {
		err = source.Close()
	}
	select {
	case r.close <- true:
	default:
	}
	return err
}
