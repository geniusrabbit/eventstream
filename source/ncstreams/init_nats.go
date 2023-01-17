//go:build nats || allsource || all
// +build nats allsource all

package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/nats"
)

func natsNewSubscriber(ctx context.Context, url string) (nc.Subscriber, error) {
	return nats.NewSubscriber(nats.WithNatsURL(url))
}

// OpenNATS new source by URLs
func OpenNATS(ctx context.Context, url string, options ...Option) (eventstream.Sourcer, error) {
	return Open(ctx, url, natsNewSubscriber, options...)
}

func init() {
	source.RegisterConnector("nats", func(ctx context.Context, config *source.Config) (eventstream.Sourcer, error) {
		return OpenNATS(ctx, config.Connect, WithDebug(config.Debug), WithFormat(config.Format))
	})
}
