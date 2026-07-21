// Package nats opens a NATS Sourcer.
package nats

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	ncnats "github.com/geniusrabbit/notificationcenter/v2/nats"
)

func newSubscriber(ctx context.Context, url string) (nc.Subscriber, error) {
	return ncnats.NewSubscriber(ncnats.WithNatsURL(url))
}

// Open connects a Sourcer.
func Open(ctx context.Context, url string, options ...source.Option) (eventstream.Sourcer, error) {
	return source.Open(ctx, url, newSubscriber, options...)
}
