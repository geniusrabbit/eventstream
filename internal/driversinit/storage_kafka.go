//go:build kafka || allstorage || all
// +build kafka allstorage all

package driversinit

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/ncstreams"
)

func init() {
	storage.RegisterConnector("kafka", func(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
		return ncstreams.OpenKafka(ctx, conf.Connect, storage.WithDebug(conf.Debug))
	})
}
