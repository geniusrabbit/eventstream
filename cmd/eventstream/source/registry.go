// Package source is the service-side source driver catalog.
// Drivers register via init() in build-tagged register_*.go files so the
// binary can include only the modules selected by build tags.
package source

import (
	"context"
	"fmt"
	"sync"

	"github.com/geniusrabbit/eventstream"
)

// Params for opening a source from a scenario item.
type Params struct {
	Connect string
	Format  string
}

// Opener opens a source for a registered driver name.
type Opener func(ctx context.Context, p Params) (eventstream.Sourcer, error)

var (
	mu      sync.RWMutex
	openers = map[string]Opener{}
)

// Register binds a driver name to an opener (called from register_*.go init).
func Register(driver string, open Opener) {
	mu.Lock()
	defer mu.Unlock()
	openers[driver] = open
}

// Open creates a source using a previously registered driver.
func Open(ctx context.Context, driver string, p Params) (eventstream.Sourcer, error) {
	mu.RLock()
	open := openers[driver]
	mu.RUnlock()
	if open == nil {
		return nil, fmt.Errorf("[cmd/source] undefined driver %q (not built into this binary)", driver)
	}
	return open(ctx, p)
}
