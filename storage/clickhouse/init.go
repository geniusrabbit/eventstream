//go:build clickhouse || allsource || all
// +build clickhouse allsource all

package clickhouse

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
)

type extraConfig struct {
	InitQuery []string `json:"init_query"`
}

// Connector of the driver
func Connector(ctx context.Context, conf *storage.Config) (eventstream.Storager, error) {
	var (
		extConf extraConfig
		err     = conf.Decode(&extConf)
	)
	if err != nil {
		return nil, err
	}
	return Open(ctx, conf.Connect,
		WithInitQuery(extConf.InitQuery),
		storage.WithDebug(conf.Debug))
}
