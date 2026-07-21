// Package natstream opens a NATS Streaming publish Storager.
package natstream

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	ncnatstream "github.com/geniusrabbit/notificationcenter/v2/natstream"
)

func newPublisher(ctx context.Context, url string) (nc.Publisher, error) {
	return ncnatstream.NewPublisher(ncnatstream.WithNatsURL(url))
}

// Open connects a publish Storager.
func Open(ctx context.Context, url string) (eventstream.Storager, error) {
	return storage.Open(ctx, url, newPublisher)
}
