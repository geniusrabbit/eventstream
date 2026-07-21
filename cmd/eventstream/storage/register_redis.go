//go:build redis || allsource || all
// +build redis allsource all

package storage

import (
	"context"
	"fmt"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage/redis"
)

func init() {
	Register("redis", func(ctx context.Context, p Params) (eventstream.Storager, error) {
		return redis.NewStorage(ctx, p.Connect)
	})
	RegisterStream("redis", func(streamName string, p StreamParams) ([]any, error) {
		if p.Key == "" {
			return nil, fmt.Errorf("stream %q: redis key required", streamName)
		}
		opts := []any{
			redis.WithID(streamName),
			redis.WithKey(p.Key),
			redis.WithIncrementor(p.Incrementor),
		}
		if p.Expiration > 0 {
			opts = append(opts, redis.WithExpiration(p.Expiration))
		}
		return opts, nil
	})
}
