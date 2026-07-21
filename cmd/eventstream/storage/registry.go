// Package storage is the service-side storage driver catalog.
// Drivers register via init() in build-tagged register_*.go files so the
// binary can include only the modules selected by build tags.
package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/geniusrabbit/eventstream"
)

// Params for opening a storage from a scenario item.
type Params struct {
	Connect   string
	Debug     bool
	InitQuery []string
	Method    string
	URL       string
}

// StreamParams for building typed Storager.Stream options from a scenario item.
type StreamParams struct {
	Debug        bool
	SQLQuery     string
	Target       string
	Fields       any
	IterateBy    string
	BufferSize   uint
	WriteTimeout uint
	Targets      []Target
	Key          string
	Expiration   time.Duration
	Incrementor  bool
	Method       string
	URL          string
	ContentType  string
}

// Target is a publish sink message shape (nats/kafka/…).
type Target struct {
	Fields map[string]any
	Where  string
}

// Opener opens a storage for a registered driver name.
type Opener func(ctx context.Context, p Params) (eventstream.Storager, error)

// StreamBuilder builds typed stream options for a store driver.
type StreamBuilder func(streamName string, p StreamParams) ([]any, error)

var (
	mu             sync.RWMutex
	openers        = map[string]Opener{}
	streamBuilders = map[string]StreamBuilder{}
)

// Register binds a driver name to an opener (called from register_*.go init).
func Register(driver string, open Opener) {
	mu.Lock()
	defer mu.Unlock()
	openers[driver] = open
}

// RegisterStream binds a driver name to a stream-options builder.
func RegisterStream(driver string, build StreamBuilder) {
	mu.Lock()
	defer mu.Unlock()
	streamBuilders[driver] = build
}

// Open creates a storage using a previously registered driver.
func Open(ctx context.Context, driver string, p Params) (eventstream.Storager, error) {
	mu.RLock()
	open := openers[driver]
	mu.RUnlock()
	if open == nil {
		return nil, fmt.Errorf("[cmd/storage] undefined driver %q (not built into this binary)", driver)
	}
	return open(ctx, p)
}

// StreamOpts builds typed stream options for a registered store driver.
func StreamOpts(driver, streamName string, p StreamParams) ([]any, error) {
	mu.RLock()
	build := streamBuilders[driver]
	mu.RUnlock()
	if build == nil {
		return nil, fmt.Errorf("[cmd/storage] undefined stream builder for driver %q (not built into this binary)", driver)
	}
	return build(streamName, p)
}
