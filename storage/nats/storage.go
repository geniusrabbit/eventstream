// Package nats opens a NATS publish Storager.
package nats

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	ncnats "github.com/geniusrabbit/notificationcenter/v2/nats"
)

func newPublisher(ctx context.Context, url string) (nc.Publisher, error) {
	return ncnats.NewPublisher(ncnats.WithNatsURL(url))
}

// Open connects a publish Storager.
func Open(ctx context.Context, url string) (eventstream.Storager, error) {
	return storage.Open(ctx, url, newPublisher)
}
