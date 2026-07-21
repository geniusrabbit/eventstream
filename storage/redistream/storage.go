// Package redistream opens a Redis Streams publish Storager.
package redistream

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/storage"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/redis"
)

func newPublisher(ctx context.Context, url string) (nc.Publisher, error) {
	return redis.NewPublisher(redis.WithRedisURL(url))
}

// Open connects a publish Storager.
func Open(ctx context.Context, url string) (eventstream.Storager, error) {
	return storage.Open(ctx, url, newPublisher)
}
