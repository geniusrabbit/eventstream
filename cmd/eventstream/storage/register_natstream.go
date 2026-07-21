//go:build natstream || allstorage || all
// +build natstream allstorage all

package storage

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage/natstream"
)

func init() {
	Register("natstream", func(ctx context.Context, p Params) (eventstream.Storager, error) {
		return natstream.Open(ctx, p.Connect)
	})
	RegisterStream("natstream", ncStreamOpts)
}
