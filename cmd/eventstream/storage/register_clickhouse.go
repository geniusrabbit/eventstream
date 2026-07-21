//go:build clickhouse || allsource || all
// +build clickhouse allsource all

package storage

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage/clickhouse"
)

func init() {
	Register("clickhouse", func(ctx context.Context, p Params) (eventstream.Storager, error) {
		opts := []clickhouse.Option{clickhouse.WithDebug(p.Debug)}
		if len(p.InitQuery) > 0 {
			opts = append(opts, clickhouse.WithInitQuery(p.InitQuery))
		}
		return clickhouse.Open(ctx, p.Connect, opts...)
	})
	RegisterStream("clickhouse", func(streamName string, p StreamParams) ([]any, error) {
		return sqlStreamOpts(streamName, p, clickhouse.WithQuery)
	})
}
