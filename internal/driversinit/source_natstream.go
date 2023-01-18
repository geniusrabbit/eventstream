//go:build natstream || allsource || all
// +build natstream allsource all

package driversinit

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	"github.com/geniusrabbit/eventstream/source/ncstreams"
)

func init() {
	source.RegisterConnector("natstream", func(ctx context.Context, config *source.Config) (eventstream.Sourcer, error) {
		return ncstreams.OpenNATStream(ctx, config.Connect,
			ncstreams.WithDebug(config.Debug),
			ncstreams.WithFormat(config.Format))
	})
}
