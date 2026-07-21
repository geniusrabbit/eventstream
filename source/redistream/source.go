// Package redistream opens a Redis Streams Sourcer.
package redistream

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/redis"
)

func newSubscriber(ctx context.Context, url string) (nc.Subscriber, error) {
	return redis.NewSubscriber(redis.WithRedisURL(url))
}

// Open connects a Sourcer.
func Open(ctx context.Context, url string, options ...source.Option) (eventstream.Sourcer, error) {
	return source.Open(ctx, url, newSubscriber, options...)
}
