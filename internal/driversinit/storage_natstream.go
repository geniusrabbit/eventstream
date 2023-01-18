//go:build natstream || allstorage || all
// +build natstream allstorage all

package driversinit

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	"github.com/geniusrabbit/eventstream/storage/ncstreams"
)

func init() {
	storage.RegisterConnector("natstream", func(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
		return ncstreams.OpenNATStream(ctx, conf.Connect, storage.WithDebug(conf.Debug))
	})
}
