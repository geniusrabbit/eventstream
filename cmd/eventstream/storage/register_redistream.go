//go:build redistream || allstorage || all
// +build redistream allstorage all

package storage

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage/redistream"
)

func init() {
	Register("redistream", func(ctx context.Context, p Params) (eventstream.Storager, error) {
		return redistream.Open(ctx, p.Connect)
	})
	RegisterStream("redistream", ncStreamOpts)
}
