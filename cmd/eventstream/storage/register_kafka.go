//go:build kafka || allstorage || all
// +build kafka allstorage all

package storage

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage/kafka"
)

func init() {
	Register("kafka", func(ctx context.Context, p Params) (eventstream.Storager, error) {
		return kafka.Open(ctx, p.Connect)
	})
	RegisterStream("kafka", ncStreamOpts)
}
