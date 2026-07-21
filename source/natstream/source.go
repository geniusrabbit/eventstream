// Package natstream opens a NATS Streaming Sourcer.
package natstream

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	ncnatstream "github.com/geniusrabbit/notificationcenter/v2/natstream"
)

func newSubscriber(ctx context.Context, url string) (nc.Subscriber, error) {
	return ncnatstream.NewSubscriber(ncnatstream.WithNatsURL(url))
}

// Open connects a Sourcer.
func Open(ctx context.Context, url string, options ...source.Option) (eventstream.Sourcer, error) {
	return source.Open(ctx, url, newSubscriber, options...)
}
