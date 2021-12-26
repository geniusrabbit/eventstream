//go:build redis || allsource || all
// +build redis allsource all

package redis

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

func init() {
	storage.RegisterConnector("redis", func(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
		return NewStorage(ctx, conf.Connect, storage.WithDebug(conf.Debug))
	})
}
