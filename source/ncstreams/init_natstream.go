//go:build natstream || allsource || all
// +build natstream allsource all

package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/natstream"
)

func natstreamNewSubscriber(ctx context.Context, url string) (nc.Subscriber, error) {
	return natstream.NewSubscriber(natstream.WithNatsURL(url))
}

// OpenNATStream new source by URLs
func OpenNATStream(ctx context.Context, url string, options ...Option) (eventstream.Sourcer, error) {
	return Open(ctx, url, natstreamNewSubscriber, options...)
}

func init() {
	source.RegisterConnector("natstream", func(ctx context.Context, config *source.Config) (eventstream.Sourcer, error) {
		return OpenNATStream(ctx, config.Connect, WithDebug(config.Debug), WithFormat(config.Format))
	})
}
