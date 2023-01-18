//go:build redis || allsource || all

package driversinit

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/redis"
)

func init() {
	storage.RegisterConnector("redis", func(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
		return redis.NewStorage(ctx, conf.Connect, storage.WithDebug(conf.Debug))
	})
}
