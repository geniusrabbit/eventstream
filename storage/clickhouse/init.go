//go:build clickhouse || allsource || all
// +build clickhouse allsource all

package clickhouse

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

func connector(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
	return Open(conf.Connect, storage.WithDebug(conf.Debug))
}

func init() {
	storage.RegisterConnector("clickhouse", connector)
}
