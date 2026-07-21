//go:build nats || allstorage || all
// +build nats allstorage all

package storage

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage/nats"
)

func init() {
	Register("nats", func(ctx context.Context, p Params) (eventstream.Storager, error) {
		return nats.Open(ctx, p.Connect)
	})
	RegisterStream("nats", ncStreamOpts)
}
