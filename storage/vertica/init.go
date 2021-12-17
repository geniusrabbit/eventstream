//go:build vertica || allsource || all
// +build vertica allsource all

package vertica

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

func connector(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
	return Open(conf.Connect, WithDebug(conf.Debug))
}

func init() {
	storage.RegisterConnector("vertica", connector)
}
