//go:build redistream || allsource || all
// +build redistream allsource all

package driversinit

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/eventstream/source/ncstreams"
)

func init() {
	source.RegisterConnector("redistream", func(ctx context.Context, config *source.Config) (eventstream.Sourcer, error) {
		return ncstreams.OpenRedis(ctx, config.Connect,
			ncstreams.WithDebug(config.Debug),
			ncstreams.WithFormat(config.Format))
	})
	source.RegisterConnector("redis", func(ctx context.Context, config *source.Config) (eventstream.Sourcer, error) {
		return ncstreams.OpenRedis(ctx, config.Connect,
			ncstreams.WithDebug(config.Debug),
			ncstreams.WithFormat(config.Format))
	})
}
