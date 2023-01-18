//go:build redistream || allsource || all
// +build redistream allsource all

package driversinit

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/ncstreams"
)

func init() {
	storage.RegisterConnector("redistream", func(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
		return ncstreams.OpenRedis(ctx, conf.Connect, storage.WithDebug(conf.Debug))
	})
}
