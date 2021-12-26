//go:build vertica || allsource || all
// +build vertica allsource all

package vertica

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	sqlstore "github.com/geniusrabbit/eventstream/storage/sql"
)

func connector(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
	return Open(conf.Connect, sqlstore.WithDebug(conf.Debug))
}

func init() {
	storage.RegisterConnector("vertica", connector)
}
