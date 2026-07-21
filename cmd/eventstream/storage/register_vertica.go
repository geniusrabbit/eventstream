//go:build vertica || allsource || all
// +build vertica allsource all

package storage

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage/vertica"
)

func init() {
	Register("vertica", func(ctx context.Context, p Params) (eventstream.Storager, error) {
		opts := []vertica.Option{vertica.WithDebug(p.Debug)}
		if len(p.InitQuery) > 0 {
			opts = append(opts, vertica.WithInitQuery(p.InitQuery))
		}
		return vertica.Open(ctx, p.Connect, opts...)
	})
	RegisterStream("vertica", func(streamName string, p StreamParams) ([]any, error) {
		return sqlStreamOpts(streamName, p, vertica.WithQuery)
	})
}
