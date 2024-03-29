//go:build redistream || allsource || all
// +build redistream allsource all

package ncstreams

import (
	"context"

	"github.com/geniusrabbit/eventstream"
	"github.com/geniusrabbit/eventstream/source"
	nc "github.com/geniusrabbit/notificationcenter/v2"
	"github.com/geniusrabbit/notificationcenter/v2/redis"
)

func redisNewSubscriber(ctx context.Context, url string) (nc.Subscriber, error) {
	return redis.NewSubscriber(redis.WithRedisURL(url))
}

// OpenRedis new source by URLs
func OpenRedis(ctx context.Context, url string, options ...Option) (eventstream.Sourcer, error) {
	return Open(ctx, url, redisNewSubscriber, options...)
}

func init() {
	source.RegisterConnector("redistream", func(ctx context.Context, config *source.Config) (eventstream.Sourcer, error) {
		return OpenRedis(ctx, config.Connect, WithDebug(config.Debug), WithFormat(config.Format))
	})
	source.RegisterConnector("redis", func(ctx context.Context, config *source.Config) (eventstream.Sourcer, error) {
		return OpenRedis(ctx, config.Connect, WithDebug(config.Debug), WithFormat(config.Format))
	})
}
