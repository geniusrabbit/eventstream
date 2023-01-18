//go:build nats || allstorage || all
// +build nats allstorage all

package driversinit

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/ncstreams"
)

func init() {
	storage.RegisterConnector("nats", func(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
		return ncstreams.OpenNATS(ctx, conf.Connect, storage.WithDebug(conf.Debug))
	})
}
